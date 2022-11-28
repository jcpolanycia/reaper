package app

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var brandedTemplateSrc = `<!DOCTYPE html>
<head>
<meta name="darkreader-lock">
<style type="text/css">
body, html {
	width: 100%;
	margin: 0;
	padding: 0;
	background-color: #000000;
	color: #ffffff;
}
.container {
	margin: 100px;
}
header {
	width: 100%;
	height: 80px;
	background-repeat: no-repeat;
	background-size: contain;
	background-image:  url('data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjwhRE9DVFlQRSBzdmcgUFVCTElDICItLy9XM0MvL0RURCBTVkcgMS4xLy9FTiIgImh0dHA6Ly93d3cudzMub3JnL0dyYXBoaWNzL1NWRy8xLjEvRFREL3N2ZzExLmR0ZCI+Cjxzdmcgd2lkdGg9IjEwMCUiIGhlaWdodD0iMTAwJSIgdmlld0JveD0iMCAwIDEwMTcgMjk4IiB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbDpzcGFjZT0icHJlc2VydmUiIHhtbG5zOnNlcmlmPSJodHRwOi8vd3d3LnNlcmlmLmNvbS8iIHN0eWxlPSJmaWxsLXJ1bGU6ZXZlbm9kZDtjbGlwLXJ1bGU6ZXZlbm9kZDtzdHJva2UtbGluZWpvaW46cm91bmQ7c3Ryb2tlLW1pdGVybGltaXQ6MjsiPgogICAgPGcgdHJhbnNmb3JtPSJtYXRyaXgoMSwwLDAsMSwtMjU2LjE1LC0xMTAuNzAxKSI+CiAgICAgICAgPGcgdHJhbnNmb3JtPSJtYXRyaXgoMC41NTI4MTEsMCwwLDAuNTUyODExLDUzNy40MjUsMTEwLjcwMSkiPgogICAgICAgICAgICA8cGF0aCBkPSJNMzg5LjI5NSwwTDMyOC4zNzYsMEMzMjguMzc2LDAgMzI4LjE3NCwwLjEwMSAzMjguMTc0LDAuMjAyTDMyOC4xNzQsMzkxLjk5QzMyOC4xNzQsMzkxLjk5IDMyOC4yNzUsMzkyLjE5MyAzMjguMzc2LDM5Mi4xOTNMMzg5LjI5NSwzOTIuMTkzQzM4OS4yOTUsMzkyLjE5MyAzODkuNDk3LDM5Mi4wOTEgMzg5LjQ5NywzOTEuOTlMMzg5LjQ5NywwLjIwMkMzODkuNDk3LDAuMjAyIDM4OS4zOTYsMCAzODkuMjk1LDBaTTQyMC4wNTcsMTc2LjE0NEM0MjIuOTU4LDE3NC4yNTUgNDI2LjA2MiwxNzIuNTY5IDQyOS4zLDE3MS4wNTFDNDM5Ljg1OCwxNjYuMTk0IDQ1MS41MjksMTY0LjAwMSA0NjMuMTMyLDE2NC40NzNDNDc0LjczNiwxNjQuOTQ2IDQ4NS4zMjcsMTY3Ljg4IDQ5NC45NDEsMTczLjI3N0M1MDUuNDk4LDE3OS4xOCA1MTQuMSwxODcuNDQ0IDUyMC43NDUsMTk4LjAwMkM1MjcuMzksMjA4LjU2IDUzMC42OTYsMjIwLjgzOCA1MzAuNjk2LDIzNC44MDNMNTMwLjY5NiwzOTIuMDI0QzUzMC42OTYsMzkyLjAyNCA1MzAuNzk3LDM5Mi4yMjYgNTMwLjg5OCwzOTIuMjI2TDU5MS43NDksMzkyLjIyNkM1OTEuNzQ5LDM5Mi4yMjYgNTkxLjk1MiwzOTIuMTI1IDU5MS45NTIsMzkyLjAyNEw1OTEuOTUyLDIzNC4yNjNDNTkxLjk1MiwyMTAuNjE4IDU4Ny4yOTcsMTg5LjYwMyA1NzcuOTg3LDE3MS4xNTJDNTY4LjY3NywxNTIuNzAxIDU1NS4xNTEsMTM4LjE5NyA1MzcuNDA4LDEyNy42MzlDNTE5LjY2NiwxMTcuMDgxIDQ5OC4yOCwxMTEuNzg1IDQ3My4yMTgsMTExLjc4NUM0NTguMTc0LDExMS43ODUgNDQzLjMzMiwxMTQuNDg0IDQyOC42MjUsMTE5Ljg0N0M0MjUuNjIzLDEyMC45MjYgNDIyLjcyMiwxMjIuMTc0IDQxOS44ODksMTIzLjUyNEM0MTkuODIxLDEyMy41MjQgNDE5Ljc4OCwxMjMuNjI1IDQxOS43ODgsMTIzLjY5Mkw0MTkuNzg4LDE3Ni4wMDlDNDE5Ljc4OCwxNzYuMDA5IDQxOS45NTYsMTc2LjI0NiA0MjAuMDkxLDE3Ni4xNzhMNDIwLjA1NywxNzYuMTQ0Wk0xMjQzLjg0LDE3MC44NDlDMTI0My44NCwxNzAuODQ5IDEyNDMuNjQsMTcwLjk1IDEyNDMuNjQsMTcxLjA1MUwxMjQzLjY0LDMxMS42MDlDMTI0My42NCwzMTkuODczIDEyNDYuMTMsMzI2LjU4NSAxMjUxLjE2LDMzMS43NDZDMTI1Ni4xOSwzMzYuOTQxIDEyNjIuOCwzMzkuNTM4IDEyNzEuMDMsMzM5LjUzOEwxMzI4LjM0LDMzOS41MzhDMTMyOC4zNCwzMzkuNTM4IDEzMjguNTQsMzM5LjYzOSAxMzI4LjU0LDMzOS43NDFMMTMyOC41NCwzOTEuOTlDMTMyOC41NCwzOTEuOTkgMTMyOC40NCwzOTIuMTkzIDEzMjguMzQsMzkyLjE5M0wxMjYwLjMsMzkyLjE5M0MxMjM1Ljk1LDM5Mi4xOTMgMTIxNi44NiwzODUuMjEgMTIwMy4wOSwzNzEuMjQ2QzExODkuMzMsMzU3LjI4MSAxMTgyLjQyLDMzOC4yOSAxMTgyLjQyLDMxNC4zMDdMMTE4Mi40Miw0OS42MTlDMTE4Mi40Miw0OS42MTkgMTE4Mi41Miw0OS40MTYgMTE4Mi42Miw0OS40MTZMMTI0My40Nyw0OS40MTZDMTI0My40Nyw0OS40MTYgMTI0My42Nyw0OS41MTcgMTI0My42Nyw0OS42MTlMMTI0My42NywxMTcuOTkyQzEyNDMuNjcsMTE3Ljk5MiAxMjQzLjc3LDExOC4xOTQgMTI0My44NywxMTguMTk0TDEzMjguOTEsMTE4LjE5NEMxMzI4LjkxLDExOC4xOTQgMTMyOS4xMSwxMTguMjk1IDEzMjkuMTEsMTE4LjM5N0wxMzI5LjExLDE3MC42NDZDMTMyOS4xMSwxNzAuNjQ2IDEzMjkuMDEsMTcwLjg0OSAxMzI4LjkxLDE3MC44NDlMMTI0My44NCwxNzAuODQ5Wk03NjEuODIyLDM5OC42MzVDNzMzLjU1NSwzOTguNjM1IDcwOC40NTksMzkyLjI2IDY4Ni42MDEsMzc5LjU3N0M2NjQuNzc3LDM2Ni44NjEgNjQ3LjQ3MywzNDkuNjkxIDYzNC43NTYsMzI4LjAwMkM2MjIuMDQsMzA2LjMxMyA2MTUuNjk4LDI4Mi4wNiA2MTUuNjk4LDI1NS4yMUM2MTUuNjk4LDIyOC4zNiA2MjIuMDczLDIwNC4xMDcgNjM0Ljc1NiwxODIuNDE4QzY0Ny40NzMsMTYwLjcyOSA2NjQuNzc3LDE0My41NiA2ODYuNjAxLDEzMC44NDNDNzA4LjQyNSwxMTguMTI3IDczMy41MjEsMTExLjc4NSA3NjEuODIyLDExMS43ODVDNzkwLjEyMiwxMTEuNzg1IDgxNC42NzksMTE4LjE2IDgzNi41MDMsMTMwLjg0M0M4NTguMzI3LDE0My41NiA4NzUuNjMxLDE2MC42NjIgODg4LjM0OCwxODIuMTQ4QzkwMS4wNjQsMjAzLjYzNSA5MDcuNDA2LDIyNy45ODkgOTA3LjQwNiwyNTUuMjFDOTA3LjQwNiwyODIuNDMxIDkwMS4wMzEsMzA2LjMxMyA4ODguMzQ4LDMyOC4wMDJDODc1LjYzMSwzNDkuNjkxIDg1OC4zMjcsMzY2Ljg2MSA4MzYuNTAzLDM3OS41NzdDODE0LjY3OSwzOTIuMjk0IDc4OS43NTEsMzk4LjYzNSA3NjEuODIyLDM5OC42MzVaTTc2MS44MjIsMzQ1Ljk4MUM3NzguNjU0LDM0NS45ODEgNzkzLjMyNywzNDEuOTY3IDgwNS44NzUsMzMzLjkwNUM4MTguNDIzLDMyNS44NDMgODI4LjEwNCwzMTUuMDE2IDgzNC44ODQsMzAxLjM4OEM4NDEuNjY0LDI4Ny43NjEgODQ1LjEwNCwyNzIuMzc5IDg0NS4xMDQsMjU1LjE3N0M4NDUuMTA0LDIzNy45NzQgODQxLjY5NywyMjMuMDMxIDgzNC44ODQsMjA5LjIzNUM4MjguMTA0LDE5NS40NzIgODE4LjQyMywxODQuNTEgODA1Ljg3NSwxNzYuNDQ4Qzc5My4zMjcsMTY4LjM4NiA3NzguNjU0LDE2NC4zNzIgNzYxLjgyMiwxNjQuMzcyQzc0NC45OSwxNjQuMzcyIDcyOS44NDUsMTY4LjM4NiA3MTcuNDk5LDE3Ni40NDhDNzA1LjE1MywxODQuNTEgNjk1LjQ3MywxOTUuMzM3IDY4OC40OSwyMDguOTY1QzY4MS41MDgsMjIyLjU5MiA2NzgsMjM3Ljk3NCA2NzgsMjU1LjE3N0M2NzgsMjcyLjM3OSA2ODEuNTA4LDI4Ny4zMjIgNjg4LjQ5LDMwMS4xMThDNjk1LjQ3MywzMTQuODgxIDcwNS4xNTMsMzI1Ljg0MyA3MTcuNDk5LDMzMy45MDVDNzI5Ljg0NSwzNDEuOTY3IDc0NC42MTksMzQ1Ljk4MSA3NjEuODIyLDM0NS45ODFaTTkzNi4wNDQsMzkxLjk5TDkzNi4wNDQsMzM5Ljc0MUM5MzYuMDQ0LDMzOS43NDEgOTM2LjE0NSwzMzkuNTM4IDkzNi4yNDYsMzM5LjUzOEwxMDY4LjIsMzM5LjUzOEMxMDc0LjY1LDMzOS41MzggMTA4MC4yOCwzMzguMTIyIDEwODUuMTQsMzM1LjI1NEMxMDg5Ljk2LDMzMi4zODcgMTA5My43NCwzMjguNjA5IDEwOTYuNCwzMjMuOTg4QzEwOTkuMSwzMTkuMzMzIDExMDAuNDIsMzE0LjUxIDExMDAuNDIsMzA5LjQ4NEMxMTAwLjQyLDMwMy43NDkgMTA5OS4xMywyOTguNzIzIDEwOTYuNjcsMjk0LjQ0QzEwOTQuMTgsMjkwLjE1NiAxMDkwLjU3LDI4Ni42NDggMTA4NS45MSwyODMuOTQ5QzEwODEuMjYsMjgxLjI1MSAxMDc1Ljg5LDI3OS45MzUgMTA2OS43OSwyNzkuOTM1TDEwMTYuMDUsMjc5LjkzNUM5OTguNTE0LDI3OS45MzUgOTgyLjkzLDI3Ny4wNjggOTY5LjMwMiwyNzEuMzM0Qzk1NS42NzUsMjY1LjU5OSA5NDQuOTQ5LDI1Ni43NjIgOTM3LjA1NSwyNDQuNzU0QzkyOS4xOTYsMjMyLjc0NSA5MjUuMjUsMjE3LjYzNCA5MjUuMjUsMTk5LjM1MUM5MjUuMjUsMTg0LjY0NSA5MjguOTI2LDE3MS4yMiA5MzYuMjQ2LDE1OS4wNDNDOTQzLjU2NiwxNDYuODY2IDk1My41MTYsMTM3LjAxNiA5NjYuMDY0LDEyOS40OTRDOTc4LjYxMiwxMjEuOTcyIDk5Mi43NDYsMTE4LjIyOCAxMDA4LjUsMTE4LjIyOEwxMTM1LjYzLDExOC4yMjhDMTEzNS42MywxMTguMjI4IDExMzUuODMsMTE4LjMyOSAxMTM1LjgzLDExOC40M0wxMTM1LjgzLDE3MC42OEMxMTM1LjgzLDE3MC42OCAxMTM1LjczLDE3MC44ODIgMTEzNS42MywxNzAuODgyTDEwMTQuNCwxNzAuODgyQzEwMDYuMTQsMTcwLjg4MiA5OTkuMzU3LDE3My40OCA5OTMuOTk0LDE3OC42NzRDOTg4LjYzLDE4My44NjkgOTg1LjkzMiwxOTAuMjEgOTg1LjkzMiwxOTcuNzMyQzk4NS45MzIsMjA1LjI1NCA5ODguNjMsMjExLjE1NyA5OTMuOTk0LDIxNi41MjFDOTk5LjM1NywyMjEuODg0IDEwMDYuNTQsMjI0LjU4MiAxMDE1LjQ4LDIyNC41ODJMMTA2Ny42LDIyNC41ODJDMTA4Ny4yOSwyMjQuNTgyIDExMDQuMTMsMjI3LjYxOCAxMTE4LjA5LDIzMy43MjNDMTEzMi4wNiwyMzkuODI5IDExNDIuNzEsMjQ4Ljk3IDExNTAuMDcsMjYxLjExM0MxMTU3LjM5LDI3My4yOSAxMTYxLjEsMjg4LjcwNSAxMTYxLjEsMzA3LjMyNUMxMTYxLjEsMzIyLjM2OSAxMTU3LjM1LDMzNi40MDEgMTE0OS44LDM0OS40ODlDMTE0Mi4yOCwzNjIuNTc3IDExMzEuOTksMzcyLjkzMiAxMTE4LjksMzgwLjY1N0MxMTA1LjgxLDM4OC4zODEgMTA5MS4wNCwzOTIuMTkzIDEwNzQuNTgsMzkyLjE5M0w5MzYuMTc4LDM5Mi4xOTNDOTM2LjE3OCwzOTIuMTkzIDkzNS45NzYsMzkyLjA5MSA5MzUuOTc2LDM5MS45OUw5MzYuMDQ0LDM5MS45OVpNMjk0LjM3NSwxMTIuMzkyTDI5NC4zNzUsNDI0LjMzOEMyOTQuMzc1LDQ2NS45NjMgMjcxLjU3Myw1MDQuMjE0IDIzNC45NzQsNTIzLjk0N0MxODguOTk5LDU0OC43NzMgMTMyLjAyNyw1MzguODU2IDk3LjExNSw0OTkuOTY0TDkwLjIzNCw0OTIuMzA3QzY4Ljg4Miw0NjguNTI2IDMzLjgzNSw0NjIuODU5IDYuMDc0LDQ3OC43MTNMMC44OCw0ODEuNjgyQzAuODgsNDgxLjY4MiAwLjU3Niw0ODEuNjgyIDAuNTc2LDQ4MS41MTNMMC41NzYsNDE3LjMyMkMwLjU3Niw0MTcuMzIyIDAuNjQ0LDQxNy4xNTQgMC43MTEsNDE3LjEyQzQ3LjQ2Myw0MDEuOTA3IDk5Ljg4MSw0MTUuMTMgMTMzLjc4MSw0NTIuOTA5TDE0MC4wODgsNDU5Ljk1OUMxNTUuNDM2LDQ3Ny4wNiAxNzkuODU3LDQ4My40MzYgMjAxLjIwOSw0NzQuOTM1QzIyMi41NjEsNDY2LjQzNSAyMzUuNjE1LDQ0Ni40MzIgMjM1LjYxNSw0MjQuMzA1TDIzNS42MTUsMzY5Ljk5OEMyMzUuNjE1LDM2OS44MjkgMjM1LjQxMywzNjkuNzI4IDIzNS4yNzgsMzY5Ljg2M0MyMzQuOTQxLDM3MC4yMzQgMjMzLjQ5LDM3MS4xMTEgMjMxLjM5OSwzNzIuMzI1QzIwNy4zODIsMzg2LjE1NSAxODAuNDMxLDM5NC4yMTcgMTUyLjc3MSwzOTYuMDcyQzE0NC45NDYsMzk2LjYxMSAxMzcuNTU5LDM5Ni45NDkgMTMxLjkyNSwzOTYuOTQ5QzEwOC4yMTIsMzk2Ljk0OSA4Ni4yMiwzOTEuMzgzIDY1Ljk0NywzODAuMjg2QzQ2LjA0NiwzNjguODE3IDMwLjE5MiwzNTIuNTI1IDE4LjMxOSwzMzEuNDc2QzYuNTQ3LDMxMC4wNTcgMC42MSwyODQuNjI0IDAuNjEsMjU1LjE0M0wwLjYxLDI0Ni41NDFDMC42MSwyMTYuNjg5IDYuNTQ3LDE5MS4yMjIgMTguMzg2LDE3MC4xNzRDMzAuNjMxLDE0OS4xMjYgNDYuNjg3LDEzMi44MzMgNjYuNTg4LDEyMS4zNjVDNzMuODQsMTE3LjE4MiA4NC40OTksMTEzLjc0MiA5NS4wOTEsMTExLjExMUMxMTkuNDc5LDEwNS4wMDUgMTQ0LjkxMiwxMDQuMDk0IDE2OS43MDQsMTA4LjE3NkwxNzAuNTgxLDEwOC4zMTFDMTkzLjc4OCwxMTIuMTU2IDIxNi4wMTcsMTIwLjM4NyAyMzYuMDU0LDEzMi43MzJDMjM2LjMyNCwxMzIuOTAxIDIzNi41NiwxMzMuMDM2IDIzNi42OTUsMTMzLjEzN0MyMzYuODI5LDEzMy4yMDUgMjM2Ljk5OCwxMzMuMTM3IDIzNi45OTgsMTMyLjk2OEwyMzYuOTk4LDExMi4zOTJDMjM2Ljk5OCwxMTIuMzkyIDIzNy4wOTksMTEyLjE5IDIzNy4yMDEsMTEyLjE5TDI5NC4xNzMsMTEyLjE5QzI5NC4xNzMsMTEyLjE5IDI5NC4zNzUsMTEyLjI5MSAyOTQuMzc1LDExMi4zOTJaTTIxMS4xOTQsMzIxLjE4OEMxOTQuNzMzLDMzNy4yNzggMTczLjcxOCwzNDUuMzA2IDE0OC4wODMsMzQ1LjMwNkMxMjIuNDQ3LDM0NS4zMDYgMTAxLjgwMywzMzcuMjc4IDg0Ljk3MiwzMjEuMTg4QzY4LjUxMSwzMDQuNzI4IDYwLjMxNCwyODIuMTI4IDYwLjMxNCwyNTMuNDIyTDYwLjMxNCwyNDguMjYyQzYwLjMxNCwyMTguNzgxIDY4LjU0NSwxOTYuMTgxIDg0Ljk3MiwxODAuNDk2QzEwMS44MDMsMTY0LjQwNiAxMjIuODUyLDE1Ni4zNzggMTQ4LjA4MywxNTYuMzc4QzE3My4zMTQsMTU2LjM3OCAxOTUuMzQsMTY0LjQwNiAyMTEuNzY3LDE4MC40OTZDMjI4LjIyOCwxOTYuMTgxIDIzNi40MjUsMjE4Ljc4MSAyMzYuNDI1LDI0OC4yNjJMMjM2LjQyNSwyNTMuNDIyQzIzNi40MjUsMjgyLjEyOCAyMjguMDI2LDMwNC43MjggMjExLjE5NCwzMjEuMTg4WiIgc3R5bGU9ImZpbGw6d2hpdGU7Ii8+CiAgICAgICAgPC9nPgogICAgICAgIDxnIHRyYW5zZm9ybT0ibWF0cml4KDAuNTc4Nzc3LDAsMCwwLjU4ODQzLC0yMDYuNzk0LC0yNDMuNjA4KSI+CiAgICAgICAgICAgIDxwYXRoIGQ9Ik0xMDE3Ljg1LDEwNDguNDNDOTk1LjU0MSwxMDQyLjY4IDk4My4yMDksMTAyOS43NiA5NzIuNDA4LDEwMTguNDRDOTY3LjI1NCwxMDEzLjAzIDk2Mi41ODMsMTAwOC4xNCA5NTguNTYyLDEwMDYuMzdDOTUzLjY0MywxMDA0LjIyIDk0NC4yOTMsMTAwNC43NCA5MzQuNTQ5LDEwMDUuMjhDOTE1Ljk3NSwxMDA2LjMyIDg5Ni4zNzUsMTAwNy40MSA4ODAuNTYyLDk5NC4xMjRDODYyLjI1MSw5NzguNzQ4IDg1Ny44NzUsOTU5LjM2IDg1My4zODMsOTM5LjQ1NkM4NDguNDgsOTE3LjczNSA4NDMuMzg4LDg5NS4xNzYgODEwLjIxMSw4ODQuMjYzQzgwMi4zNDgsODgxLjcwNyA3OTguMDQ2LDg3My4yNiA4MDAuNjAzLDg2NS4zOTdDODAyLjI4OSw4NjAuMjEgODA2LjUzOSw4NTYuNTczIDgxMS40OTMsODU1LjQzMkM4NjIuNDIxLDg0MS4xMDcgODkzLjUzMyw4MjEuOTQyIDkxNC45NDQsNzk5LjgxMkM5MzYuMzksNzc3LjY0NiA5NDguODc2LDc1MS40NDQgOTYxLjg2Nyw3MjQuMTA3Qzk3OS4yNDMsNjg3LjUzIDEwMDQuMTYsNjYzLjkwNSAxMDMxLjI5LDY1Mi4wNzNDMTA1MS41OSw2NDMuMjE4IDEwNzMuMDIsNjQwLjk4NyAxMDkzLjQ3LDY0NC44ODJDMTExMy44Nyw2NDguNzY3IDExMzMuMjIsNjU4LjcxMiAxMTQ5LjQxLDY3NC4yMThDMTE3Mi4zLDY5Ni4xMzYgMTE4OC44NSw3MjkuMjI1IDExOTIuNzksNzcxLjkyQzExOTYuMjcsODA5LjY5MyAxMTkyLjE5LDg1Mi4zMzcgMTE4MS44NSw4OTIuNTExQzExNzEuNjUsOTMyLjA5MyAxMTU1LjMxLDk2OS41MjcgMTEzNC4wMiw5OTcuOTAyQzExMTEuNDksMTAyNy45NCAxMDgzLjI4LDEwNDguMTYgMTA1MC41LDEwNTEuMTlDMTAzOS45NywxMDUyLjE2IDEwMjkuMDcsMTA1MS4zMiAxMDE3Ljg1LDEwNDguNDNaTTk5NC4wODcsOTk3LjgxNEMxMDAyLjE1LDEwMDYuMjYgMTAxMS4zNSwxMDE1LjkxIDEwMjUuMjMsMTAxOS40OUMxMDMzLjA1LDEwMjEuNSAxMDQwLjU5LDEwMjIuMDkgMTA0Ny44MSwxMDIxLjQyQzEwNzEuMywxMDE5LjI2IDEwOTIuNDcsMTAwMy41IDExMTAuMTIsOTc5Ljk3M0MxMTI5LjAxLDk1NC43ODUgMTE0My42NSw5MjEuMDcgMTE1Mi45LDg4NS4xMjhDMTE2Mi40Myw4NDguMTQxIDExNjYuMiw4MDkuMDU1IDExNjMuMDIsNzc0LjYxNUMxMTU5LjgsNzM5LjY5NCAxMTQ2LjcxLDcxMy4wNTMgMTEyOC42Nyw2OTUuNzhDMTExNi43Nyw2ODQuMzgyIDExMDIuNjksNjc3LjEgMTA4Ny45Niw2NzQuMjk1QzEwNzMuMjgsNjcxLjQ5OSAxMDU3Ljg3LDY3My4xMTMgMTA0My4yNCw2NzkuNDkzQzEwMjIuNCw2ODguNTgzIDEwMDIuOTYsNzA3LjM1OCA5ODguOTM2LDczNi44OEM5NzQuODYxLDc2Ni41IDk2MS4zMjYsNzk0Ljg5NSA5MzYuMzg4LDgyMC42N0M5MTYuODUzLDg0MC44NjIgODkwLjc5Miw4NTguNzkzIDg1Mi44MDksODczLjMwOEM4NzIuODY1LDg4OS45NjIgODc3Ljc4MSw5MTEuNzE2IDg4Mi41NjEsOTMyLjg5NEM4ODUuOTE5LDk0Ny43NzIgODg5LjE5LDk2Mi4yNjQgODk5Ljc3OSw5NzEuMTU3QzkwNi41NzgsOTc2Ljg2OCA5MjAuMDk3LDk3Ni4xMTUgOTMyLjkwOCw5NzUuNDAxQzk0Ni4yNjYsOTc0LjY1NyA5NTkuMDgzLDk3My45NDIgOTcwLjUxNCw5NzguOTU0Qzk4MC4xMSw5ODMuMTY0IDk4Ni43NTUsOTkwLjEyOSA5OTQuMDg3LDk5Ny44MTRaIiBzdHlsZT0iZmlsbDp3aGl0ZTtmaWxsLXJ1bGU6bm9uemVybzsiLz4KICAgICAgICA8L2c+CiAgICAgICAgPGcgdHJhbnNmb3JtPSJtYXRyaXgoMC41Nzg3NzcsMCwwLDAuNTg4NDMsLTIwNi43OTQsLTI0My42MDgpIj4KICAgICAgICAgICAgPHBhdGggZD0iTTEwNjQuMSw3NDkuNDI2QzEwNzQuOTQsNzQ5LjQyNiAxMDgzLjczLDc1OC4yMTcgMTA4My43Myw3NjkuMDU5QzEwODMuNzMsNzc5LjkwMiAxMDc0Ljk0LDc4OC42OTMgMTA2NC4xLDc4OC42OTNDMTA1My4yNSw3ODguNjkzIDEwNDQuNDYsNzc5LjkwMiAxMDQ0LjQ2LDc2OS4wNTlDMTA0NC40Niw3NTguMjE3IDEwNTMuMjUsNzQ5LjQyNiAxMDY0LjEsNzQ5LjQyNloiIHN0eWxlPSJmaWxsOndoaXRlOyIvPgogICAgICAgIDwvZz4KICAgICAgICA8ZyB0cmFuc2Zvcm09Im1hdHJpeCgwLjU3ODc3NywwLDAsMC41ODg0MywtMjA2Ljc5NCwtMjQzLjYwOCkiPgogICAgICAgICAgICA8cGF0aCBkPSJNMTEyMC45Niw3NDkuNDI2QzExMzEuOCw3NDkuNDI2IDExNDAuNTksNzU4LjIxNyAxMTQwLjU5LDc2OS4wNTlDMTE0MC41OSw3NzkuOTAyIDExMzEuOCw3ODguNjkzIDExMjAuOTYsNzg4LjY5M0MxMTEwLjEyLDc4OC42OTMgMTEwMS4zMiw3NzkuOTAyIDExMDEuMzIsNzY5LjA1OUMxMTAxLjMyLDc1OC4yMTcgMTExMC4xMiw3NDkuNDI2IDExMjAuOTYsNzQ5LjQyNloiIHN0eWxlPSJmaWxsOndoaXRlOyIvPgogICAgICAgIDwvZz4KICAgIDwvZz4KPC9zdmc+Cg==')
}
main {
	padding-top: 60px;
}
a {
	color: #492CFB;
}
</style>
</head>
<body>
	<div class="container">
		<header></header>
		<main>{{.Msg}}</main>
	</div>
</body>
</html>`

var brandedTemplate = template.Must(template.New("branded").Parse(brandedTemplateSrc))

type templateInput struct {
	Msg      template.HTML
	Hostname string
}

func (a *App) handleLocalRequest(request *http.Request) *http.Response {

	caDownload := base64.RawURLEncoding.EncodeToString(a.userSettings.Get().CACert)
	caDownloadURL := fmt.Sprintf("data:application/octet-stream;base64,%s", caDownload)

	return a.createReaperMessageResponse(request, fmt.Sprintf(`
<p>Welcome to Reaper!</p>

<p>Download the CA certificate for your proxy here: <a download="reaper-ca.crt" href="%s">Download CA Certificate</a>.</p>

<p>For further information and configuration, please use the Reaper desktop application.</p>

`, caDownloadURL))
}

func (a *App) createRawMessageResponse(req *http.Request, data io.Reader, contentType string) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    req,
		Body:       io.NopCloser(data),
		Header: map[string][]string{
			"Content-Type": {contentType},
		},
	}
}

func (a *App) createReaperMessageResponse(req *http.Request, msg string) *http.Response {
	buf := bytes.NewBuffer(nil)
	if err := brandedTemplate.Execute(buf, templateInput{
		Msg:      template.HTML(msg),
		Hostname: a.userSettings.Get().ProxyHost,
	}); err != nil {
		return a.createRawMessageResponse(req, strings.NewReader(msg), "text/plain")
	}
	return a.createRawMessageResponse(req, buf, "text/html")
}
