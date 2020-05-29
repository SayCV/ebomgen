# ebomgen

## 1. Introduction

ebomgen is a tool to auto generate bom from EDA design file, it support Kicad, Eagle, Orcad, Altium or Mentor Graphics etc.

## 2. Getting started

### 2.1. Install

```BATCH
git clone github.com/saycv/ebomgen
cd ebomgen && make install
```

### 2.2. Usage

```BATCH
echo Process PADSLogic ASCII file
ebomgen -t padslogic -i test/pads/SCH/ex1.txt -o test/pads/BOM
echo Process PADSPCB ASCII file
ebomgen -t padspcb -i test/pads/PCB/ex2.asc -o test/pads/BOM

echo Process AltiumPCB ASCII file
ebomgen -t altiumpcb -i test/altium/PCB/ex4.PcbDoc -o test/altium/BOM
```

## 3. TODO

- [x] Process Altium file
- [ ] Process Eagle file
- [ ] Process Kicad file
- [ ] Process OrCAD file
- [x] Process PADSLogic or PADSPCB file
- [x] Export CSV file
- [ ] Export Exls file
- [x] Sorted Components
- [x] Sorted References
- [ ] Calc MTBF based PCP
- [ ] More humane
- [ ] More smart

## 4. Example

```CSV
Item,References,Quantity,Value,Footprint,Description
1,"C1,C2",2,"???"," DCAP1","Capacitor"
2,"C3,C4,C5,C6,C7,C8,C9,C10,C11,C12,C13,C14,C15,C16,C17,C18,C19,C20,C21,C22",20,".01UF"," DCAP1","Capacitor"
3,"R1,R2",2,"1K"," R1/8W","Resistor"
4,"Y1",1,"XTAL1"," XTAL1","Crystal"
5,"U12,U13,U14,U15,U16,U17,U18,U19,U20,U21",10,"6167"," DIP20","IC"
6,"U1",1,"74139"," DIP16","IC"
7,"U9,U11",2,"7400"," DIP14","IC"
8,"U2,U6",2,"7404"," DIP14","IC"
9,"U3",1,"7402"," DIP14","IC"
10,"U7",1,"7440"," DIP14","IC"
11,"U8",1,"7420"," DIP14","IC"
12,"U4",1,"7474"," DIP14","IC"
13,"U5",1,"7432"," DIP14","IC"
14,"P1",1,"CON\26P\ED"," CON\26P\ED","Connector"
```
