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
```

## 3. TODO

- [ ] Process Altium file
- [ ] Process Eagle file
- [ ] Process Kicad file
- [ ] Process OrCAD file
- [x] Process PADSLogic or PADSPCB file
- [ ] Calc MTBF
