
## 🟠 Partie 2 : Persistance des données & Communication (50 min)

*Objectif : Rendre l'application interactive. Vous allez devoir réceptionner les données envoyées par le navigateur, les valider, les sauvegarder dans une base de données locale (SQLite) pour garder une trace, et enfin les transmettre au serveur distant.*

### Étape 2.1 : Les Modèles de données
Pour gérer les échanges, vous devez définir deux structures dans le fichier `models/submission.go`.

**1. La structure de soumission (`Submission`)**
Elle représente le dessin envoyé par l'utilisateur.
*   `ID` (uint, clé primaire).
*   `Name` (string) : Obligatoire (`binding:"required"`).
*   `Grid` ([][]string) : Reçoit la grille brute depuis le JSON. **Attention :** SQL ne gère pas ce type. Utilisez le tag `gorm:"-"` pour l'ignorer en base.
*   `GridData` (string) : Servira à stocker la grille convertie en texte (JSON stringifié) dans la BDD.
*   `CreatedAt` (time.Time).

**2. La structure de réponse API (`APIResponse`)**
Le frontend (le fichier HTML/JS fourni) s'attend à recevoir une réponse JSON standardisée pour afficher les messages dans la zone "status".
Définissez une structure `APIResponse` contenant :
*   `Success` (bool) : Indique si l'opération a réussi.
*   `Message` (string) : Le texte explicatif qui s'affichera sur l'écran de l'utilisateur.
*   *N'oubliez pas les tags json correspondants (`json:"success"`, etc.).*

### Étape 2.2 : Connexion à la Base de Données
Utilisez un singleton (variable globale) pour gérer la connexion.

1.  Créez le fichier `database/db.go`.
2.  Déclarez une variable globale `DB` de type `*gorm.DB`.
3.  Implémentez une fonction `Connect()` qui :
    *   Ouvre une connexion SQLite (fichier `pixel.db`).
    *   Utilise `DB.AutoMigrate(...)` pour créer la table `Submission`.
    *   Gère les erreurs de connexion.
4.  **Intégration dans le main :** Allez immédiatement dans votre fichier `main.go` et ajoutez l'appel à `database.Connect()` **au tout début** de la fonction `main()`.

> ⚠️ **Attention :** Si vous oubliez d'appeler `database.Connect()` dans le `main`, la variable `DB` restera vide (`nil`). Votre programme **crashera** (runtime error / panic) dès que vous tenterez de sauvegarder une grille à l'étape suivante.

### Étape 2.3 : Envoi au serveur distant (Service)
Dans `services/api_proxy.go`, ajoutez la fonction pour contacter l'API du professeur.

```go
// PostGridToRemote envoie les données au serveur distant
// payload correspond à votre structure Submission
func PostGridToRemote(payload any) ([]byte, int, error) {
    // 1. Convertir le payload en JSON (Marshal)
    // 2. Faire une requête POST sur ServerAPI + "/api/submit"
    // 3. Retourner le body de la réponse et le status code
}
```

### Étape 2.4 : La Soumission (Contrôleur)
C'est le cœur du projet. Dans `controllers/pixel_controller.go`, créez la fonction `SubmitProxyGrid` (Route `POST /proxy/submit`).

**Algorithme à implémenter :**

1.  **Binding :** Récupérez le JSON dans la structure `Submission`.
    *   Si le binding échoue (ex: nom manquant), renvoyez une erreur 400 en utilisant votre structure `APIResponse` (Success: false, Message: "Erreur...").
2.  **Préparation :** Convertissez le champ `Grid` (tableau) en `string` (via `json.Marshal`) et stockez-le dans `GridData`.
3.  **Sauvegarde :** Enregistrez la soumission en local avec `database.DB.Create`.
    *   En cas d'erreur SQL, renvoyez une 500 avec `APIResponse`.
4.  **Envoi Distant :** Appelez votre service `PostGridToRemote`.
5.  **Réponse Final :**
    *   Si l'envoi distant échoue, prévenez l'utilisateur mais confirmez la sauvegarde locale via un `APIResponse`.
    *   Sinon, renvoyez directement la réponse brute reçue du serveur distant.

### Étape 2.5 : L'Historique local
Permettez à l'utilisateur de voir ses anciens dessins.

1.  **API (`GetLocalHistory`)** :
    *   Récupérez les **10 dernières soumissions** depuis la BDD (`created_at desc`).
    *   Retournez la liste en JSON.
    *   Route : `GET /proxy/history`.
2.  **HTML (`RenderHistory`)** :
    *   Affichez simplement le template `history.html`.
    *   Route : `GET /history`.
    




Remplacer history.html par le fichier suivant: 

```
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Historique Pixel Art</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; background: #f0f2f5; padding: 20px; }

        h1 { text-align: center; color: #333; }

        .nav-link { display: block; text-align: center; margin-bottom: 30px; text-decoration: none; color: #3498db; font-weight: bold; }
        .nav-link:hover { text-decoration: underline; }

        /* Conteneur des cartes */
        .gallery {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 20px;
            max-width: 1200px;
            margin: 0 auto;
        }

        /* Une carte individuelle */
        .card {
            background: white;
            border-radius: 10px;
            padding: 15px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            display: flex;
            flex-direction: column;
            align-items: center;
        }

        .card-header {
            width: 100%;
            display: flex;
            justify-content: space-between;
            margin-bottom: 10px;
            font-size: 0.9em;
            color: #666;
            border-bottom: 1px solid #eee;
            padding-bottom: 5px;
        }
        .user-name { font-weight: bold; color: #333; }

        /* La mini grille */
        .mini-grid {
            display: grid;
            grid-template-columns: repeat(15, 10px); /* Cellules de 10px */
            gap: 0;
            border: 2px solid #333;
            background: #ccc;
        }

        .mini-cell {
            width: 10px;
            height: 10px;
            background-color: white;
        }
    </style>
</head>
<body>

<h1>🏛️ Galerie des Œuvres</h1>
<a href="/" class="nav-link">← Retour au jeu</a>

<div class="gallery" id="gallery-container">
    <!-- Les cartes seront injectées ici par JS -->
    <p style="text-align:center; width:100%;">Chargement des données...</p>
</div>

<script>
    // URL de l'API locale qu'on a créée dans l'étape précédente
    const API_HISTORY_URL = "/proxy/history";

    async function loadHistory() {
        try {
            const res = await fetch(API_HISTORY_URL);
            if (!res.ok) throw new Error("Erreur réseau");

            const submissions = await res.json();
            renderGallery(submissions);
        } catch (e) {
            document.getElementById('gallery-container').innerHTML =
                `<p style="color:red; text-align:center;">Impossible de charger l'historique : ${e.message}</p>`;
        }
    }

    function renderGallery(submissions) {
        const container = document.getElementById('gallery-container');
        container.innerHTML = '';

        if (submissions.length === 0) {
            container.innerHTML = '<p>Aucune donnée en base.</p>';
            return;
        }

        submissions.forEach(sub => {
            // 1. Création de la carte
            const card = document.createElement('div');
            card.className = 'card';

            // 2. Parsing de la date
            const dateObj = new Date(sub.created_at);
            const dateStr = dateObj.toLocaleDateString() + ' ' + dateObj.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});

            // 3. En-tête de la carte
            card.innerHTML = `
                <div class="card-header">
                    <span class="user-name">👤 ${sub.name}</span>
                    <span>${dateStr}</span>
                </div>
            `;

            // 4. Génération de la mini-grille
            // IMPORTANT : Dans la BDD, la grille est stockée en string ("[[...]]"), il faut la parser
            let gridData = [];
            try {
                // Si sub.GridData est vide, on met une grille vide
                gridData = sub.GridData ? JSON.parse(sub.GridData) : [];
            } catch(e) {
                console.error("Erreur parsing grille", e);
            }

            const gridDiv = document.createElement('div');
            gridDiv.className = 'mini-grid';

            // Dessiner les 15x15 cellules
            // Si la grille récupérée n'est pas complète, on gère l'affichage vide
            for(let r=0; r<15; r++) {
                for(let c=0; c<15; c++) {
                    const cell = document.createElement('div');
                    cell.className = 'mini-cell';

                    // On vérifie si la donnée existe à ces coordonnées
                    if(gridData[r] && gridData[r][c]) {
                        cell.style.backgroundColor = gridData[r][c];
                    }
                    gridDiv.appendChild(cell);
                }
            }

            card.appendChild(gridDiv);
            container.appendChild(card);
        });
    }

    // Lancer le chargement au démarrage
    loadHistory();
</script>

</body>
</html>
```

---

### ✅ Vérification
1.  Lancez le serveur.
    *   *Si le serveur crash immédiatement ou au moment de l'envoi, vérifiez que vous avez bien fait l'étape 2.2 point 4.*
2.  Essayez d'envoyer un dessin **sans mettre de nom** : le message d'erreur doit s'afficher grâce à votre `APIResponse`.
3.  Envoyez un dessin valide : vous devez recevoir le succès.
4.  Vérifiez l'onglet "Historique local".





### Étape 2.5 : Modification...

Modifiez la fonction GetProxyState  et FetchStateFromRemote comme suit :

```
func GetProxyState(c *gin.Context) {
body, status, _ := services.FetchStateFromRemote()
c.Data(status, "application/json", body)
}



func FetchStateFromRemote() ([]byte, int, error) {
    req, _ := http.NewRequest(http.MethodGet, ServerAPI+"/api/state", nil)
    resp, _ := httpClient.Do(req)
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    return body, resp.StatusCode, nil
}
```

Après avoir modifié le code, faites signe à l'un des étudiants animant le cours pour qu'il puisse procéder à la démo.

Important : Attendez les instructions avant de continuer...

