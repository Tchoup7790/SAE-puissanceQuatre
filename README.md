# Puissance Quatre

## Branch Baptiste

Bienvenue dans le projet Puissance Quatre en Golang !
Ce projet vise à recréer le  jeu Puissance 4 en utilisant le langage de programmation Golang. 
Deux joueurs peuvent s'affronter via un serveur dédié.

## Comment jouer

1. **Déplacement dans le répertoire :** Utilisez la commande suivante pour aller dans le bon répertoire.
    ```bash
    cd ./client
    ```
   
2. **Instalation des dépendances :** Utilisez la commande suivante pour intaller les dépendances du projet.
    ```bash
    go install
    ```

3. **Build du projet :** Utilisez la commande suivante pour construire le projet.
    ```bash
    go build
    ```

4. **Lancement du jeu :** Après avoir construit le projet, lancez le jeu avec la commande suivante.
    ```bash
    ./client
    ```

5. **Jouer :** Suivez les instructions à l'écran pour jouer votre tour. Le jeu prend en charge deux joueurs qui 
alternent pour placer leurs jetons sur le plateau de jeu.


## Options de lancemen
### -addr
Changer l'address ip et le port du **client**
```bash
./client -addr=ip:port
```
### -debug
Activer le mod debug du **client**
```bash
./client -debug
```


## Issues

Si vous rencontrer un problème ou que vous avez des requêtes suggestions à nous faire, n'hésitez pas à ouvrir une issue sur le repository.
