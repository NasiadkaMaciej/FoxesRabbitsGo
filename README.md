# Foxes and Rabbits Simulation

| Main simulation view | Population chart |
| ---- | ----------- |
| ![](https://nasiadka.pl/projects/FoxesRabbitsGo/main.png) | ![](https://nasiadka.pl/projects/FoxesRabbitsGo/chart.png) |


## Project Description
This project presents an ecosystem simulation in which foxes and rabbits coexist in an environment with regrowable grass.

## Features
- Interactive simulation with graphical visualization
- Foxes hunt rabbits for energy
- Rabbits eat grass and flee from foxes. They may "stumble" while fleeing
- Animals can reproduce when appropriate conditions are met
  - Fox reproduction conditions:
    - Fox energy must be higher than reproduction cost
    - Another fox must be nearby
    - A specific number of turns must have passed since previous reproduction
  - Rabbit reproduction conditions:
    - Rabbit energy must be higher than reproduction cost
    - Another rabbit must be nearby
    - A specific number of turns must have passed since previous reproduction
- Grass regrows over time
  - When grass is eaten, a countdown to its regrowth begins
  - Grass regrows to a specific amount in one cell

## Installation
```bash
go mod init
```

## Running
```bash
go run main.go
```

## Controls
- Left mouse button: Add a rabbit
- Right mouse button: Add a fox
- You can "draw" animals

## Configuration
Simulation parameters can be changed in the config.go file.