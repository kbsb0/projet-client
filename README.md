# 🎓 Sujet de TP : Pixel Proxy Client
**Technologies :** Go, Gin Framework, HTML/JS, HTTP Client.

## 🎯 Objectif
L'objectif de ce TP est de développer un **client API intermédiaire (Proxy)**. Votre application devra :
1.  Afficher une interface web permettant de dessiner une grille de pixels.
2.  Interroger un serveur distant pour connaître le modèle à dessiner.
3.  Sauvegarder vos dessins dans une base de données locale (SQLite).
4.  Transmettre vos créations au serveur distant.

**Concepts abordés :** Routing avec Gin, architecture MVC, Rendu HTML, appels HTTP (Client), JSON Binding.


## 🟢 Partie 1 : Initialisation & Routing (30 min)

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


### Étape 1.3 : Structure et Rendu HTML
Nous allons maintenant structurer le projet et afficher l'interface graphique.

---

## 📂 Ressources préliminaires

Avant de commencer le développement Go, vous devez mettre en place l'interface utilisateur.
Créez un dossier nommé `templates` à la racine de votre projet et ajoutez-y les deux fichiers suivants.


index.hmtl :
```html
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>{{ .title }}</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; display: flex; flex-direction: column; align-items: center; background: #f0f2f5; }
        .top-info { display: flex; gap: 20px; align-items: center; margin: 15px 0; }
        .timer { font-weight: bold; color: #e74c3c; font-size: 1.5em; }
        input { padding: 8px; border: 1px solid #ccc; border-radius: 4px; }
        .grid-container { display: grid; grid-template-columns: repeat(15, 30px); gap: 1px; background: #ccc; border: 5px solid #333; user-select: none; }
        .cell { width: 30px; height: 30px; background: white; cursor: pointer; }
        .cell.hint { background-color: #dcdcdc; }
        .palette { display: flex; gap: 10px; margin: 15px 0; }
        .color-choice { width: 35px; height: 35px; border-radius: 50%; cursor: pointer; border: 3px solid transparent; box-shadow: 0 2px 5px rgba(0,0,0,0.2); }
        .color-choice.active { transform: scale(1.2); border-color: #333; }
        button { padding: 10px 20px; cursor: pointer; font-size: 16px; font-weight: bold; border:none; border-radius: 4px; color: white; margin-top: 15px; }
        .btn-send { background-color: #2ed573; }
        .btn-clear { background-color: #ff4757; }
        #message { font-weight: bold; margin-top: 15px; height: 20px; text-align: center;}
    </style>
</head>
<body>

<h1>{{ .title }}</h1>

<div class="top-info">
    <div class="timer" id="client-timer">--:--</div>
    <input type="text" id="username" placeholder="Votre Prénom" maxlength="15">
</div>

<div class="palette">
    <div class="color-choice active" style="background-color: #3498db;" data-color="#3498db"></div>
    <div class="color-choice" style="background-color: #e74c3c;" data-color="#e74c3c"></div>
    <div class="color-choice" style="background-color: #f1c40f;" data-color="#f1c40f"></div>
    <div class="color-choice" style="background-color: #2ecc71;" data-color="#2ecc71"></div>
    <div class="color-choice" style="background-color: #9b59b6;" data-color="#9b59b6"></div>
    <div class="color-choice" style="background-color: #34495e;" data-color="#34495e"></div>
</div>

<div class="grid-container" id="grid"></div>

<div class="controls">
    <button class="btn-clear" onclick="clearGrid()">Tout Effacer</button>
    <button class="btn-send" onclick="sendGrid()">ENVOYER</button>
</div>

<div id="message"></div>

<script>
    // Configuration
    const API_STATE_URL = "/proxy/state";
    const API_SUBMIT_URL = "/proxy/submit";

    let currentColor = '#3498db';
    let currentModelId = -1;
    let localTimeLeft = 0;
    let syncInProgress = false;

    // --- INITIALISATION UI ---

    // Gestion de la Palette de couleurs
    document.querySelectorAll('.color-choice').forEach(c => {
        c.addEventListener('click', () => {
            document.querySelectorAll('.color-choice').forEach(x => x.classList.remove('active'));
            c.classList.add('active');
            currentColor = c.dataset.color;
        });
    });

    // Initialisation d'une grille vide au démarrage
    drawGrid(Array(15).fill().map(() => Array(15).fill(0)));


    // --- PARTIE A DECOMMENTER - ETAPE 1.3 (Proxy State) ---
    /*
    async function syncState() {
        if(syncInProgress) return;
        syncInProgress = true;
        try {
            const res = await fetch(API_STATE_URL);
            if (!res.ok) throw new Error("Erreur proxy ou Serveur distant");

            const data = await res.json();

            localTimeLeft = data.timeLeft || 0;
            updateTimerDisplay();

            if (data.targetGrid && currentModelId !== data.currentModel) {
                currentModelId = data.currentModel;
                drawGrid(data.targetGrid);
            }
        } catch(e) {
            console.warn("Sync error:", e);
        } finally {
            syncInProgress = false;
        }
    }

    function countdown() {
        if (localTimeLeft > 0) {
            localTimeLeft--;
            updateTimerDisplay();
        }
        if (localTimeLeft <= 0) {
            syncState();
        }
    }

    // Lancer la boucle de synchro
    setInterval(countdown, 1000);
    syncState();
    */

    // --- FONCTIONS UTILITAIRES ---

    function updateTimerDisplay() {
        let t = Math.floor(localTimeLeft);
        if (t < 0) t = 0;
        let m = Math.floor(t / 60);
        let s = t % 60;
        document.getElementById('client-timer').innerText = (m<10?'0':'')+m + ":" + (s<10?'0':'')+s;
    }

    function drawGrid(model) {
        const container = document.getElementById('grid');
        container.innerHTML = '';
        for(let r=0; r<15; r++) {
            for(let c=0; c<15; c++) {
                const d = document.createElement('div');
                d.className = 'cell' + (model && model[r] && model[r][c] === 1 ? ' hint' : '');
                d.onclick = function() {
                    const previousColor = this.style.backgroundColor;
                    if (this.style.backgroundColor === currentColor ||
                        (this.style.backgroundColor === 'rgb(' + hexToRgb(currentColor) + ')')) {
                        this.style.backgroundColor = "";
                    } else {
                        this.style.backgroundColor = currentColor;
                    }
                };
                container.appendChild(d);
            }
        }
    }

    function hexToRgb(hex) {
        var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
        return result ?
            parseInt(result[1], 16) + ", " + parseInt(result[2], 16) + ", " + parseInt(result[3], 16)
            : null;
    }

    function clearGrid() {
        document.querySelectorAll('.cell').forEach(c => c.style.backgroundColor = '');
    }

    // --- PARTIE A DECOMMENTER - ETAPE 2.3 (Envoi du dessin) ---
    async function sendGrid() {
        const name = document.getElementById('username').value.trim();
        if(!name) { alert("Merci d'entrer votre prénom !"); return; }

        // Préparation des données
        let grid = [];
        for(let r=0; r<15; r++) {
            let row = [];
            for(let c=0; c<15; c++) {
                let cell = document.getElementById('grid').children[r*15 + c];
                row.push(cell.style.backgroundColor || "");
            }
            grid.push(row);
        }

        /*
           DECOMMENTER CI-DESSOUS POUR ENVOYER AU SERVEUR
        */

        /*
        const msgDiv = document.getElementById('message');
        msgDiv.innerText = "Envoi en cours...";

        try {
            const res = await fetch(API_SUBMIT_URL, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({ grid: grid, name: name })
            });

            const data = await res.json();
            if (!res.ok) throw new Error(data.message || "Erreur serveur");

            msgDiv.innerText = data.message || "Dessin envoyé !";
            msgDiv.style.color = data.success ? "#2ed573" : "#ff4757";

            // Force une mise à jour immédiate
            if(typeof syncState === "function") await syncState();

        } catch(e) {
            msgDiv.innerText = "Erreur : " + e.message;
            msgDiv.style.color = "#ff4757";
            console.error(e);
        }
        */

        // Alert temporaire tant que le code ci-dessus est commenté
        alert("L'envoi vers le serveur n'est pas encore activé ! (Voir Etape 2.3)");
    }
</script>
</body>
</html>
```


history.hmtl :
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

Fonction FetchStateFromRemote

But : Récupérer l’état du serveur distant /api/state et renvoyer au contrôleur les informations nécessaires pour répondre au client.

| Valeur   | Type     | Signification                                                         |
| -------- | -------- | --------------------------------------------------------------------- |
| `[]byte` | `body`   | Contenu brut de la réponse HTTP (JSON du serveur distant)             |
| `int`    | `status` | Code HTTP de la réponse distante (ex : 200, 404, 503)                 |
| `error`  | `err`    | Erreur éventuelle lors de la requête (connexion impossible, timeout…) |


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
