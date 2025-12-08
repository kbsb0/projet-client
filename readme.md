C'est noté ! Voici le sujet de TP réécrit sous forme de **consignes pédagogiques**. Je ne donne pas le code final du Go, mais les squelettes (`boilerplate`) et les instructions pour que tu puisses l'implémenter toi-même.

Je fournis par contre les fichiers HTML complets (comme demandé) pour que tu n'aies pas à faire de front-end.

---

# 🔵 TP Partie 3 : Authentification & Sécurité (Sujet)

**Durée estimée :** 1h15
**Objectif :** Transformer notre application de Pixel Art "naïve" (où n'importe qui peut mettre n'importe quel nom) en une application sécurisée.
Nous allons implémenter :
1.  Une base de données d'utilisateurs.
2.  Un système d'inscription (hashage de mot de passe).
3.  Un système de login (JWT stocké dans un Cookie).
4.  Un middleware pour protéger les routes.
5.  L'utilisation de l'identité connectée pour signer les dessins.

---

## 📂 Pré-requis : Mise en place des fichiers

### 1. Installation des dépendances
Ouvrez votre terminal et installez les paquets pour gérer les mots de passe et les tokens :
```bash
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/golang-jwt/jwt/v5
```

### 2. Les fichiers HTML (Templates)
Dans votre dossier `templates/`, assurez-vous d'avoir les 4 fichiers suivants.
*Note : `index.html` et `history.html` sont ceux que tu as fournis (avec les modifications pour gérer la redirection si non connecté), voici les deux nouveaux :*

#### A. `templates/register.html` (Nouveau)
```html
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Inscription</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background: #f0f2f5; }
        .card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; width: 300px; }
        input { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ccc; box-sizing: border-box;}
        button { background: #27ae60; color: white; border: none; padding: 10px; width: 100%; cursor: pointer; }
        .error { color: red; margin-bottom: 10px; font-size: 0.9em; }
        a { display: block; margin-top: 10px; color: #3498db; text-decoration: none; }
    </style>
</head>
<body>
    <div class="card">
        <h2>📝 Inscription</h2>
        {{ if .error }}<div class="error">{{ .error }}</div>{{ end }}
        <form action="/register" method="POST">
            <input type="text" name="username" placeholder="Choisissez un pseudo" required>
            <input type="password" name="password" placeholder="Mot de passe" required>
            <button type="submit">S'inscrire</button>
        </form>
        <a href="/login">Déjà un compte ? Se connecter</a>
    </div>
</body>
</html>
```

#### B. `templates/login.html` (Nouveau)
```html
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Connexion</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background: #f0f2f5; }
        .card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; width: 300px; }
        input { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ccc; box-sizing: border-box;}
        button { background: #3498db; color: white; border: none; padding: 10px; width: 100%; cursor: pointer; }
        .error { color: red; margin-bottom: 10px; font-size: 0.9em; }
        a { display: block; margin-top: 10px; color: #666; text-decoration: none; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="card">
        <h2>🔐 Connexion</h2>
        {{ if .error }}<div class="error">{{ .error }}</div>{{ end }}
        <form action="/login" method="POST">
            <input type="text" name="username" placeholder="Pseudo" required>
            <input type="password" name="password" placeholder="Mot de passe" required>
            <button type="submit">Se connecter</button>
        </form>
        <a href="/register">Créer un compte</a>
    </div>
</body>
</html>
```

---

## 📝 Étape 1 : Le Modèle Utilisateur (10 min)

Nous devons stocker les utilisateurs en base de données.

1.  Créez le fichier `models/user.go`.
2.  Définissez une structure `User` qui hérite de `gorm.Model`.
3.  Ajoutez les champs :
    *   `Username` (string) : doit être unique (indice Gorm `uniqueIndex`).
    *   `Password` (string) : stockera le **hash**, pas le clair !
4.  Dans `database/database.go`, ajoutez `&models.User{}` dans la fonction `AutoMigrate` pour créer la table au démarrage.

---

## 🔐 Étape 2 : Inscription et Connexion (30 min)

Créez le fichier `controllers/auth.go`. Nous allons gérer la logique d'entrée/sortie.

### 2.1 Inscription (`Register`)
Implémentez la fonction qui reçoit le formulaire POST.

*   Récupérez `username` et `password` via `c.PostForm(...)`.
*   **Sécurité :** Utilisez `bcrypt.GenerateFromPassword` pour hasher le mot de passe.
*   Créez l'utilisateur en BDD.
*   En cas d'erreur (ex: pseudo déjà pris), réaffichez la template `register.html` avec un message d'erreur.
*   En cas de succès, redirigez vers `/login`.

### 2.2 Connexion (`Login`)
Implémentez la fonction qui vérifie les identifiants.

*   Cherchez l'utilisateur dans la BDD par son `username`.
*   **Vérification :** Utilisez `bcrypt.CompareHashAndPassword` pour comparer le hash stocké et le mot de passe reçu.
*   **Création du Token :**
    *   Utilisez la librairie `jwt-go` (v5).
    *   Créez des `claims` (données) contenant le `username` et une date d'expiration (`exp`).
    *   Signez le token avec une clé secrète (ex: une constante globale).
*   **Stockage :** Placez ce token dans un **Cookie** via `c.SetCookie(...)`.
    *   *Astuce :* Mettez `HttpOnly` à `true` pour empêcher le vol de cookie par JavaScript.

```go
// Squelette de controllers/auth.go
var jwtKey = []byte("ma_super_cle_secrete")

func Register(c *gin.Context) {
    // TODO: Récupérer form -> Hasher password -> Sauver User -> Redirect Login
}

func Login(c *gin.Context) {
    // TODO: Trouver User -> Comparer Hash -> Créer JWT -> SetCookie -> Redirect Home
}

func Logout(c *gin.Context) {
    // TODO: Écraser le cookie avec une durée de vie négative -> Redirect Login
}
```

---

## 👮 Étape 3 : Middleware d'Authentification (20 min)

Nous devons intercepter les requêtes pour vérifier si l'utilisateur est connecté.

1.  Créez `middlewares/auth.go`.
2.  Implémentez `AuthMiddleware() gin.HandlerFunc`.

**Logique à implémenter :**
1.  Récupérez le cookie nommé "auth_token" (`c.Cookie(...)`).
2.  S'il n'y a pas de cookie : redirigez vers `/login` et avortez la requête (`c.Abort()`).
3.  Parsez le token avec `jwt.Parse`.
4.  Vérifiez si le token est valide. Si non -> Redirect login.
5.  **Crucial :** Extrayez le `username` des claims du token et stockez-le dans le contexte Gin :
    ```go
    c.Set("username", claims["username"])
    ```
    *Cela permettra aux contrôleurs suivants de savoir QUI est connecté.*
6.  Laissez passer la requête avec `c.Next()`.

---

## 🚀 Étape 4 : Adaptation des Routes et Contrôleurs (15 min)

### 4.1 Mise à jour de `main.go`
Organisez vos routes.
*   Les routes `/login`, `/register` doivent être publiques.
*   Les routes `/`, `/history` et `/proxy/...` doivent être dans un **Groupe** qui utilise votre middleware.

### 4.2 Modification de `SubmitProxyGrid`
Dans `controllers/pixel.go`, la fonction `SubmitProxyGrid` reçoit actuellement le nom de l'utilisateur via le JSON (`req.Name`). **C'est une faille de sécurité**, n'importe qui peut se faire passer pour un autre.

1.  Modifiez la fonction pour ignorer le champ `name` du JSON.
2.  Récupérez le vrai nom de l'utilisateur connecté via le contexte :
    ```go
    username, exists := c.Get("username")
    ```
3.  Utilisez ce `username` pour créer l'objet `Submission`.

### 4.3 Modification de `RenderHome`
L'index.html a besoin d'afficher le pseudo de l'utilisateur (`{{ .username }}`).
Modifiez `RenderHome` pour récupérer le username du contexte (`c.Get`) et le passer à `c.HTML`.

---

## 🧪 Étape 5 : Test (Bonus)

Lancez votre serveur (`go run .`).

1.  Tentez d'aller sur `http://localhost:8081/`. Vous devriez être redirigé vers le Login.
2.  Créez un compte "Toto".
3.  Connectez-vous.
4.  Sur la page de dessin, vérifiez que votre pseudo s'affiche en haut.
5.  Dessinez et envoyez.
6.  Vérifiez dans l'historique que c'est bien "Toto" qui a signé l'œuvre.
7.  Essayez de modifier le code JS dans la console du navigateur pour envoyer un autre nom : le serveur doit l'ignorer et utiliser "Toto" quand même.