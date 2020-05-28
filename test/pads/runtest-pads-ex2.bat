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
ebomgen -t padspcb -i PCB/ex2.asc -o BOM/ >BOM/ex2.convert.log 2>&1 && echo ok || echo fail
pause
