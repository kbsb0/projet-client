Voici une proposition de structuration pour ton TP **"Pixel Proxy Client"**. Le projet est découpé en **3 grandes parties** progressives pour tenir dans le créneau de 2h.

Je fournis d'abord les **Ressources (HTML)** à donner aux étudiants dès le début, puis le déroulé des exercices.

---

# 📂 Ressources à fournir aux étudiants (Dès le début)

Les étudiants doivent créer un dossier `templates` à la racine et y placer ces deux fichiers. Cela leur évite de perdre du temps sur le HTML/JS.

<details>
<summary>📄 <b>templates/index.html</b> (Cliquer pour voir)</summary>

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
    <style>
        body { font-family: sans-serif; text-align: center; }
        .grid { display: grid; grid-template-columns: repeat(10, 30px); gap: 2px; justify-content: center; margin: 20px auto; }
        .cell { width: 30px; height: 30px; border: 1px solid #ccc; cursor: pointer; }
        .controls { margin: 20px; }
    </style>
</head>
<body>
    <h1>{{ .title }}</h1>
    <div class="controls">
        <input type="text" id="username" placeholder="Votre nom" />
        <input type="color" id="colorPicker" value="#000000">
        <button onclick="submitGrid()">Envoyer le dessin</button>
        <a href="/history"><button>Voir l'historique local</button></a>
    </div>
    <div id="grid" class="grid"></div>
    <div id="status"></div>

    <script>
        const grid = document.getElementById('grid');
        let gridData = Array(10).fill().map(() => Array(10).fill("#ffffff"));

        // Init Grid UI
        for(let i=0; i<10; i++) {
            for(let j=0; j<10; j++) {
                let cell = document.createElement('div');
                cell.className = 'cell';
                cell.onclick = () => {
                    let color = document.getElementById('colorPicker').value;
                    cell.style.backgroundColor = color;
                    gridData[i][j] = color;
                };
                grid.appendChild(cell);
            }
        }

        async function submitGrid() {
            const name = document.getElementById('username').value;
            const res = await fetch('/proxy/submit', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: name, grid: gridData })
            });
            const result = await res.json();
            document.getElementById('status').innerText = result.message || "Envoyé !";
        }
    </script>
</body>
</html>
```
</details>

<details>
<summary>📄 <b>templates/history.html</b> (Cliquer pour voir)</summary>

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .title }}</h1>
    <a href="/">Retour</a>
    <ul id="list"></ul>
    <script>
        fetch('/proxy/history')
            .then(res => res.json())
            .then(data => {
                const list = document.getElementById('list');
                data.forEach(item => {
                    let li = document.createElement('li');
                    li.innerText = `${item.created_at} - ${item.name}`;
                    list.appendChild(li);
                });
            });
    </script>
</body>
</html>
```
</details>

---

# 🎓 Sujet du TP : Pixel Proxy Client (2h)

**Objectif :** Créer un client API intermédiaire (Proxy) qui permet de dessiner une grille de pixels, de la sauvegarder dans une base de données locale (SQLite), et de l'envoyer vers un serveur distant.

**Concepts GIN abordés :** Routing, Grouping, HTML Rendering, JSON Binding, Middleware, Testing.

---

## 🟢 Partie 1 : Initialisation & Routing (40 min)
*Objectif : Mettre en place le serveur, servir les fichiers HTML et créer une première route API simple.*

### Étape 1.1 : Setup du projet
1. Créez un dossier `ari2-client`.
2. Initialisez le module Go : `go mod init ari2-client`.
3. Installez Gin et GORM (pour plus tard) :
   ```bash
   go get -u github.com/gin-gonic/gin
   go get -u gorm.io/gorm
   go get -u gorm.io/driver/sqlite
   ```

### Étape 1.2 : Serveur Web & HTML
Dans le fichier `main.go`, configurez un serveur Gin de base qui :
1. Charge les templates HTML situés dans `templates/*`.
2. Définit une route `GET /` qui appelle une fonction contrôleur `RenderHome`.
3. Cette fonction (à créer dans `controllers/pixel_controller.go`) doit rendre le fichier `index.html` avec un titre "Pixel Challenge Pro".

> 💡 **Aide :** Utilisez `r.LoadHTMLGlob` et `c.HTML`.

### Étape 1.3 : La couche Service & Proxy simple
Nous voulons récupérer l'état du serveur distant.
1. Créez le fichier `services/api_proxy.go`.
2. Implémentez la fonction `FetchStateFromRemote()` (Code fourni ci-dessous à compléter/analyser).
3. Créez une route API dans `main.go` sous un groupe `/proxy` : `GET /proxy/state`.
4. Créez le contrôleur `GetProxyState` qui appelle le service et renvoie le JSON brut au client.

**Code à utiliser pour `services/api_proxy.go` :**
```go
package services

import (
	"io"
	"net/http"
	"time"
)

const ServerAPI = "http://localhost:8080" // Serveur distant fictif

var httpClient = &http.Client{ Timeout: 5 * time.Second }

func FetchStateFromRemote() ([]byte, int, error) {
	// TODO: Faire un GET sur ServerAPI + "/api/state"
    // TODO: Lire le body et le retourner avec le status code
    // (Voir le code complet fourni si besoin d'aide)
}
```

---

## 🟠 Partie 2 : Binding JSON & Base de données (50 min)
*Objectif : Gérer les données entrantes (POST), valider le JSON, sauvegarder en local et envoyer au distant.*

### Étape 2.1 : Le Modèle de données
Créez `models/submission.go`.
Définissez la structure `Submission` qui servira à la fois pour le JSON (reçu du frontend) et GORM (BDD).
*   Attention : Le champ `Grid` est un tableau 2D (`[][]string`), difficile à stocker tel quel en SQL.
*   Ajoutez un champ `GridData` (string) pour la BDD et utilisez le tag `gorm:"-"` sur `Grid` pour l'ignorer en base.

### Étape 2.2 : Connexion BDD
Créez `database/db.go`.
1. Créez une variable globale `DB`.
2. Implémentez `Connect()` qui ouvre `sqlite.db` et lance `AutoMigrate(&models.Submission{})`.
3. Appelez `Connect()` au début du `main.go`.

### Étape 2.3 : Soumission de grille (Le cœur du projet)
Dans `controllers/pixel_controller.go`, implémentez la fonction `SubmitProxyGrid`.
C'est une route `POST /proxy/submit`.

**La logique à implémenter :**
1. **Binding :** Récupérer le JSON envoyé par le client dans la struct `Submission`.
    *   *Défi :* Si le nom est vide, renvoyer une erreur 400 (`binding:"required"`).
2. **Préparation BDD :** Convertir `Submission.Grid` (tableau) en JSON string pour le mettre dans `Submission.GridData`.
3. **Sauvegarde :** Utiliser `database.DB.Create(...)` pour sauver en local.
4. **Envoi distant :** Appeler une nouvelle fonction service `PostGridToRemote` (à créer dans `services/`) qui envoie les données à l'API distante.
5. **Réponse :** Renvoyer le résultat final au client.

### Étape 2.4 : Historique local
1. Créez la route `GET /proxy/history` et son contrôleur `GetLocalHistory`.
2. Elle doit renvoyer les 10 dernières soumissions enregistrées en base (JSON).
3. Ajoutez la route HTML `GET /history` qui affiche `history.html`.

---

## 🔵 Partie 3 : Middleware & Testing (30 min)
*Objectif : Fiabiliser l'application avec des logs structurés et des tests unitaires.*

### Étape 3.1 : Middleware Custom (Logger)
Créez `middlewares/custom.go`.
Implémentez un middleware `RequestLogger` qui :
1. Génère un UUID unique pour chaque requête.
2. Ajoute cet ID dans le Header de réponse `X-Request-ID`.
3. Loggue dans la console : `[REQ ID] METHOD PATH | STATUS | LATENCY`.
4. Appliquez ce middleware globalement dans `main.go` avec `r.Use()`.

### Étape 3.2 : Test Unitaire (Validation)
On veut s'assurer que l'API rejette bien une grille sans nom d'utilisateur.
Créez le fichier `main_test.go`.

**Exercice :**
1. Créez une fonction `SetupRouter()` qui retourne un `*gin.Engine` configuré juste avec la route `/proxy/submit`.
2. Écrivez `TestSubmitProxyGrid_Validation` :
    *   Créez une requête POST avec un JSON invalide (champ `name` manquant).
    *   Utilisez `httptest.NewRecorder()`.
    *   Assert que le Code de retour est bien `400` (BadRequest).

---

## 🚀 Pour aller plus loin (Bonus si temps disponible)

*   **Gestion d'erreur avancée :** Si le serveur distant est éteint, la sauvegarde locale doit quand même fonctionner (mode dégradé). Vérifiez que votre code le gère.
*   **Affichage de la grille dans l'historique :** Modifier le template `history.html` et le contrôleur pour décoder le JSON stocké en base et afficher un petit aperçu des couleurs.

---

### Résumé des fichiers à produire par l'étudiant :

1.  `main.go` (Point d'entrée, routes)
2.  `database/db.go` (Connexion SQLite)
3.  `models/submission.go` (Structures)
4.  `controllers/pixel_controller.go` (Logique métier)
5.  `services/api_proxy.go` (Appels HTTP sortants)
6.  `middlewares/custom.go` (Log & UUID)
7.  `main_test.go` (Test)