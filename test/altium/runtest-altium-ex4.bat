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
ebomgen -t altiumpcb -i PCB/ex4.PcbDoc -o BOM/ >BOM/ex4.convert.log 2>&1 && echo ok || echo fail && pause
