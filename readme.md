
# 😈 TP Partie 4 : Le "Cheat Mode" (Introduction à la Concurrence)

**Objectif :** Vous en avez marre de cliquer case par case ? Nous allons implémenter un bouton "TRICHER" qui remplit automatiquement la grille et l'envoie **5 fois simultanément** au serveur distant pour maximiser vos points (ou spammer le serveur).

**Concepts abordés :**
*   Goroutines (`go func()`)
*   Synchronisation (`sync.WaitGroup`)
*   Décodage JSON côté Go (`json.Unmarshal`)

---

## 🛠 Étape 1 : Modification de l'interface (HTML)

Nous devons ajouter un bouton discret pour déclencher la triche.

Dans le fichier `templates/index.html`, ajoutez ce bouton dans la `div` `.controls`, juste à côté du bouton "ENVOYER" :

```html
<!-- Dans templates/index.html -->
<div class="controls">
  <button class="btn-clear" onclick="clearGrid()">Tout Effacer</button>
  <button class="btn-send" onclick="sendGrid()">ENVOYER</button>
  
  <!-- NOUVEAU BOUTON -->
  <button style="background-color: #8e44ad;" onclick="cheat()">😈 TRICHER</button>
</div>
```

Ensuite, ajoutez la fonction JavaScript `cheat()` dans la balise `<script>` (tout en bas) pour appeler votre future route Go :

```javascript
// Dans templates/index.html (script)

async function cheat() {
    const msgDiv = document.getElementById('message');
    msgDiv.innerText = "Lancement des robots...";
    msgDiv.style.color = "#8e44ad";

    try {
        // On appelle la route Go qui va gérer la concurrence
        const res = await fetch("/proxy/cheat", { 
            method: 'POST' 
        });
        
        const data = await res.json();
        msgDiv.innerText = data.message;
        
    } catch (e) {
        msgDiv.innerText = "Erreur triche : " + e.message;
    }
}
```

---

## 🛣 Étape 2 : La Route de la Triche

Dans `main.go`, ajoutez la route correspondante dans votre groupe protégé (`protected`).

```go
// Dans main.go, sous le groupe "api"
api.POST("/cheat", controllers.CheatHandler)
```

---

## 🧠 Étape 3 : Le Contrôleur (Le cœur du sujet)

C'est ici que vous allez coder la logique concurrente. Créez (ou complétez) le contrôleur.

**Algorithme à implémenter dans `CheatHandler` :**

1.  **Récupérer la solution :** Appelez `services.FetchStateFromRemote()` pour obtenir l'état actuel du jeu (qui contient la grille cible).
2.  **Décoder le JSON :** Le serveur renvoie des `bytes`. Vous devez les transformer en une structure Go pour lire la grille cible (`TargetGrid`).
3.  **Préparer l'envoi :** Construisez un objet `Submission` parfait (ou rempli d'une couleur unique si vous préférez simplifier).
4.  **Lancer les Goroutines :**
    *   Utilisez un `sync.WaitGroup`.
    *   Lancez une boucle de 5 itérations.
    *   À chaque itération, lancez une `go func()` qui envoie la grille via `services.PostGridToRemote()`.
5.  **Attendre :** Bloquez l'exécution tant que les 5 requêtes ne sont pas finies avec `wg.Wait()`.
6.  **Répondre :** Envoyez un JSON au client disant "5 soumissions envoyées !".

### Aide pour le code

Voici les structures nécessaires pour décoder la réponse du serveur distant (à mettre dans `models/game.go` ou dans le contrôleur) :

```go
// Structure pour lire la réponse de /api/state
type RemoteState struct {
    TimeLeft     float64   `json:"timeLeft"`
    CurrentModel int       `json:"currentModel"`
    TargetGrid   [][]int   `json:"targetGrid"` // 0 ou 1
}
```

Voici le squelette de la fonction à compléter dans `controllers/pixel.go` :

```go
func CheatHandler(c *gin.Context) {
    // 1. Récupérer l'état distant
    body, _, err := services.FetchStateFromRemote()
    if err != nil {
        c.JSON(500, gin.H{"error": "Impossible de lire l'état distant"})
        return
    }

    // 2. Décoder le JSON pour avoir la grille cible
    var state models.RemoteState
    json.Unmarshal(body, &state)

    // Astuce : La TargetGrid du serveur contient des 0 et des 1.
    // Pour tricher, transformons les '1' en une couleur (ex: bleu "#3498db")
    // et les '0' en blanc "".
    var cheatGrid [][]string
    for _, row := range state.TargetGrid {
        var colorRow []string
        for _, cell := range row {
            if cell == 1 {
                colorRow = append(colorRow, "#3498db") // On force le bleu
            } else {
                colorRow = append(colorRow, "")
            }
        }
        cheatGrid = append(cheatGrid, colorRow)
    }

    username, _ := c.Get("username")

    // Préparation de la structure d'envoi
    submission := models.Submission{
        Name: username.(string), // Cast interface{} vers string
        Grid: cheatGrid,
    }

    // --- DÉBUT DE LA ZONE DE CONCURRENCE ---
    
    // TODO: Initialiser le WaitGroup

    // TODO: Faire une boucle de 0 à 5
        // TODO: wg.Add(1)
        // TODO: Lancer la goroutine (go func() { ... })
            // TODO: defer wg.Done()
            // TODO: Appeler services.PostGridToRemote(submission)
    
    // TODO: Attendre la fin des goroutines (wg.Wait())
    
    // --- FIN DE LA ZONE DE CONCURRENCE ---

    c.JSON(200, gin.H{
        "message": "💥 5 Grilles envoyées en parallèle !",
    })
}
```

### 💡 Question bonus
*Pourquoi utilise-t-on `wg.Wait()` avant d'envoyer la réponse `c.JSON` ? Que se passerait-il si on l'enlevait ?*