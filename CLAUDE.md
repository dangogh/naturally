# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go package implementing natural sort ordering for string slices. Numbers within strings are compared numerically rather than lexicographically (e.g., "A2" < "A11"). Inspired by Perl's Sort::Naturally.

## Commands

- **Run tests:** `go test ./...`
- **Run a single test:** `go test -run TestSortsA`
- **Run benchmarks:** `go test -bench .`

## Architecture

Single-package library with no dependencies beyond the standard library. The core type is `StringSlice` (a `sort.StringSlice` alias) which implements `sort.Interface`. The `Less` method splits strings into alternating non-digit/digit segments and compares them appropriately (string comparison for text, numeric comparison for numbers).
