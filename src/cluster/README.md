# Cluster Module
  - Cluster represents a set of Microcontrollers running instances of GoFire

  - Clusters are responsible for:
    - Determining and set Master/Slave state of this microcontroller
    - If this microcontroller is in master mode, managing de/registration of other microcontrollers, Id generation and handling commands for local components. 
    - If this microcontroller is in slave mode, reporting to master and handling commands for local components. 

