---

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

---

### ✅ Vérification
1.  Lancez le serveur.
    *   *Si le serveur crash immédiatement ou au moment de l'envoi, vérifiez que vous avez bien fait l'étape 2.2 point 4.*
2.  Essayez d'envoyer un dessin **sans mettre de nom** : le message d'erreur doit s'afficher grâce à votre `APIResponse`.
3.  Envoyez un dessin valide : vous devez recevoir le succès.
4.  Vérifiez l'onglet "Historique local".