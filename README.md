# 🎓 Sujet de TP : Pixel Proxy Client
**Technologies :** Go, Gin Framework, HTML/JS, HTTP Client.

## 🎯 Objectif
L'objectif de ce TP est de développer un **client API intermédiaire (Proxy)**. Votre application devra :
1.  Afficher une interface web permettant de dessiner une grille de pixels.
2.  Interroger un serveur distant pour connaître le modèle à dessiner.
3.  Sauvegarder vos dessins dans une base de données locale (SQLite).
4.  Transmettre vos créations au serveur distant.

**Concepts abordés :** Routing avec Gin, architecture MVC, Rendu HTML, appels HTTP (Client), JSON Binding.


## 🟢 Partie 1 : Initialisation & Routing (40 min)

*Objectif : Mettre en place le serveur web, vérifier son fonctionnement et servir les fichiers HTML.*

### Étape 1.1 : Configuration du projet
1.  Créez un dossier nommé `ari2-client`.
2.  Initialisez le module Go via votre terminal :
    ```bash
    go mod init ari2-client
    ```
3.  Installez les dépendances nécessaires (Gin et GORM) :
    ```bash
    go get -u github.com/gin-gonic/gin
    go get -u gorm.io/gorm
    go get -u gorm.io/driver/sqlite
    ```

### Étape 1.2 : Vérification de l'environnement (Hello World)
Avant d'intégrer les templates, nous allons créer un serveur minimaliste pour s'assurer que tout fonctionne correctement.

1.  Créez un fichier `main.go` à la racine.
2.  Insérez le code suivant :
    ```go
    package main

    import (
        "github.com/gin-gonic/gin"
    )

    func main() {
        // Création du routeur avec les middlewares par défaut (logger + recovery)
        r := gin.Default()

        // Route de test
        r.GET("/", func(c *gin.Context) {
            c.String(200, "Hello world depuis Gin !")
        })

        // Lancement du serveur sur le port 8081
        r.Run(":8081")
    }
    ```
3.  Lancez le serveur : `go run main.go`
4.  Ouvrez votre navigateur à l’adresse : [http://localhost:8081](http://localhost:8081). Vous devriez voir le message de bienvenue.


### Étape 1.3 : Structure MVC et Rendu HTML
Nous allons maintenant structurer le projet et afficher l'interface graphique.

---

## 📂 Ressources préliminaires

Avant de commencer le développement Go, vous devez mettre en place l'interface utilisateur.
Créez un dossier nommé `templates` à la racine de votre projet et ajoutez-y les deux fichiers suivants.

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
        
        // Note pour l'étudiant : Ce script pourra être enrichi plus tard pour récupérer l'état initial.
    </script>
</body>
</html>
```

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

1.  Créez un dossier `controllers` et, à l'intérieur, un fichier `pixel_controller.go`.
2.  Dans ce contrôleur, créez une fonction `RenderHome` qui prend en paramètre le contexte Gin (`*gin.Context`). Cette fonction doit afficher le template `index.html` avec le titre "Pixel Challenge Pro".
3.  Modifiez le fichier `main.go` pour :
    *   Charger les templates HTML situés dans `templates/*` (utilisez `r.LoadHTMLGlob`).
    *   Remplacer la route de test précédente par l'appel à `RenderHome`.

> 💡 **Aide :** La méthode `c.HTML(http.StatusOK, "nom_du_fichier", data)` permet de rendre une vue.

Une fois cette étape terminée, en rafraîchissant la page [http://localhost:8081](http://localhost:8081), vous devriez voir apparaître la grille de dessin.

### Étape 1.4 : Service Proxy & Récupération de l'état
Vous avez la grille, mais vous ne savez pas encore quel dessin réaliser. Cette information est détenue par le serveur distant (API du professeur).

Nous allons créer une route "Proxy" : votre navigateur demandera l'info à votre serveur Go, qui la demandera au serveur distant.

**Architecture de la requête :**
`Navigateur` -> `GET /proxy/state` (Votre serveur) -> `GET /api/state` (Serveur Distant)

1.  Créez le fichier `services/api_proxy.go`.
2.  Implémentez la fonction `FetchStateFromRemote()` en complétant le code ci-dessous.
3.  Dans `main.go`, créez un groupe de routes `/proxy` et ajoutez la route `GET /proxy/state`.
4.  Créez un contrôleur `GetProxyState` (dans `pixel_controller.go`) qui appelle votre service et retourne le JSON brut au client.

**Code squelette pour `services/api_proxy.go` :**
```go
package services

import (
	"io"
	"net/http"
	"time"
)

// Le serveur distant tourne sur le port 8080
const ServerAPI = "http://localhost:8080" 

var httpClient = &http.Client{ Timeout: 5 * time.Second }

func FetchStateFromRemote() ([]byte, int, error) {
	// TODO: Faire une requête GET sur ServerAPI + "/api/state"
    // TODO: Lire le corps de la réponse (Body)
    // TODO: Retourner les données (byte array), le code HTTP et l'erreur éventuelle
    
    // Indice : utilisez http.NewRequest, httpClient.Do, et io.ReadAll
    return nil, 0, nil
}
```

**Vérification :**
Pour tester que votre proxy fonctionne, assurez-vous que le serveur distant est lancé, puis accédez à [http://localhost:8081/](http://localhost:8081/). Vous devriez voir apparaître de nouveau la grille, mais cette fois-ci avec le dessin à réaliser.
