---
title: '10 Forms of Bash Shell Parameter Expansion'
date: 2023-10-15T10:00:00+06:00
tags:
  - Bash
  - 100DaysToOffload
toc: yes
---

If I look at my search history with the word "bash" in it, the most frequently searched phrases turn out to be like "trim suffix bash" and "set bash variable if empty".

It seems I write Bash scripts frequently enough to need these, but not frequently enough to remember these simple Bash shell parameter expansion forms.

In this blog post I am going to keep a list of 10 forms of Bash shell parameter expansion, that I hope will save me a visit to google.com.

In Bash, the basic form of parameter expansion is `${var}`. But there's more. The official documentation has a detailed list and description of how else you can expand a parameter: https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html

## `${var:-fallback}`: Use `fallback`, If `var` Is Unset

``` bash
A=
B=cookies
C=cakes
echo ${A:-keyboardcat} # ⇒ keyboardcat
echo ${B:-keyboardcat} # ⇒ cookies
echo ${A:-$C}          # ⇒ cakes
```

## `${var:+substitute}`: Use `substitute`, If `var` Is Set

``` bash
A=
B=cookies
C=cakes
echo ${A:+keyboardcat} # ⇒ 
echo ${B:+keyboardcat} # ⇒ keyboard
echo ${C:+$B}          # ⇒ cookies
```

## `${#var}`: Length of `var`

``` bash
A="the quick brown fox"
echo ${#A} # ⇒ 19
```

## `${var:start}`: Substring of `var` Starting from `start` 

``` bash
A="the quick brown fox"
echo ${A:10} # ⇒ brown fox
```

## `${var:start:length}`: Substring of `var` of Size `length` Starting from `start`

``` bash
A="the quick brown fox"
echo ${A:4:5} # ⇒ quick
```

## `${var#prefix}`: Trim `prefix` From `var`

``` bash
FILENAME="/home/hjr265/not-important-do-not-open.txt"
DIR="/home/hjr265/"
echo ${FILENAME#/home/hjr265/} # ⇒ not-important-do-not-open.txt
echo ${FILENAME#$DIR}          # ⇒ not-important-do-not-open.txt
```

## `${var%suffix}`: Trim `suffix` From `var`

``` bash
FILENAME="not-important-do-not-open.txt"
EXT=".txt"
echo ${FILENAME%.txt} # ⇒ not-important-do-not-open
echo ${FILENAME%$EXT} # ⇒ not-important-do-not-open
```

## `${var/pattern/substitute}`: Replace `pattern` in `var` with `substitute`

``` bash
A="the slow brown fox"
B="fox"
C="cow"
echo ${A/slow/quick} # => the quick brown fox
echo ${A/$B/$C}      # => the slow brown cow
```

## `${var^}`: Change Letter Case in `var`

``` bash
A="the quick brown fox"
B="THE SLOW BROWN COW"
echo ${A^}  # => The quick brown fox
echo ${A^^} # => THE QUICK BROWN FOX
echo ${B,}  # => tHE SLOW BROWN COW
echo ${B,,} # => the slow brown cow
```

## `${var@Q}`: Quoted String

```  bash
A="-w -s"
echo ${A@Q} # => '-w -s'
```
