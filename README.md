# Projet RSA 2023

## Sujet

En partant du couple client-serveur TCP, développer une application simple de tchat
d’assistance multi-niveaux, avec les trois niveaux suivants :

- le niveau 1 correspond à un robot automatique, qui traite les requêtes simples (nous
demandons un fonctionnement très basique du robot avec reconnaissance de
certains mots clés),
- le niveau 2 envoie la requête à des techniciens pour la résolution d’incidents courants
(si le niveau 1 ne trouve pas de solution),
- le niveau 3 envoie la requête à des experts du domaine pour la résolution d’incidents
plus complexes (si aucun des techniciens de niveau 2 ne trouvent de solution).

Le nombre de techniciens et d’experts est variable, ceux-ci sont localisés sur des sites
différents et sont contactés en parallèle (pour un même niveau).

## Interprétation du sujet

Nous avons choisis de créer un protocole avec des ["paquets" au format JSON](./hub/readme.md) bien définis.

Nous avons aussi décider de faire 4 programmes, un programme pour le serveur central appelé "Hub", un programme pour les clients, un programme pour les agents, et un programme pour les robots. Les robots utilisent les mêmes routes que les agents mais sont programmés pour répondre qu'à un certains nombre de client (voir [ici](./bot/readme.md)).

## Mise en place d'un réseau

- [Hub](./hub/install.md)
- [Client](./client/install.md)
- [Agent](./agent/install.md)
- [Bot](./bot/install.md)
