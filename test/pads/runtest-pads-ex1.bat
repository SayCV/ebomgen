@echo off
echo,
echo,=====
echo,SPDX-License-Identifier: (GPL-2.0+ OR MIT):
echo,
echo,!!! THIS IS NOT GUARANTEED TO WORK !!!
echo,
echo,Copyright (c) 2018-2020, SayCV
echo,=====
echo,

set "PATH=../../bin;%PATH%"

if not exist "BOM" mkdir BOM
ebomgen -t padslogic -i SCH/ex1.txt -o BOM/ >BOM/ex1.convert.log 2>&1 && echo ok || echo fail && pause

ebomgen bomcost -i BOM/ex1_BOM.csv -o BOM/ex1_BOM.bomcost.csv >BOM/ex1.bomcost.log 2>&1 && echo ok || echo fail && pause

call ebomgen bommtbf ^
		    -i BOM/ex1_BOM.csv ^
		    -o BOM/ex1_BOM.mtbf.30.csv ^
            -q "C1" ^
            -c "GB" ^
            -e "60" ^
            -d "0.7" ^
		    >>BOM/ex1.bommtbf.30.log 2>&1 && echo ok || echo fail && pause

call ebomgen bommtbf ^
		    -i BOM/ex1_BOM.csv ^
		    -o BOM/ex1_BOM.mtbf.40.csv ^
            -q "C1" ^
            -c "GF1" ^
            -e "70" ^
            -d "0.7" ^
		    >>BOM/ex1.bommtbf.40.log 2>&1 && echo ok || echo fail && pause

call ebomgen bommtbf ^
		    -i BOM/ex1_BOM.csv ^
		    -o BOM/ex1_BOM.mtbf.55.csv ^
            -q "C1" ^
            -c "GM1" ^
            -e "85" ^
            -d "0.7" ^
		    >>BOM/ex1.bommtbf.55.log 2>&1 && echo ok || echo fail && pause
