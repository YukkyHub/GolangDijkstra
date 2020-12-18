# GolangProj

Projet réalisé par Jean Vignat, Reynald Lambolez, Bryan Djafer et Mathis Chapuis en 3TCA.

Pour exécuter le programme, saisir la commande suivante :

```
go run src/main.go INPUT_FILE [OUTPUT_FILE]
```

Si aucun fichier de sortie n'est donné, les résultats seront écrits dans **dijsktra-output.txt**.

## Performances

Commme discuté pendant le vocal, voici les performances en faisant varier le nombre de workers :

* Avec 1 worker
  * Environ 6.7s d'exécution
* Avec 2 workers
  * Environ 4 secondes
* Avec 10 workers
  * Environ 2.7s
