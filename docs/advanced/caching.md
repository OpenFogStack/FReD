---
layout: default
title: NaSe Cache
parent: Advanced Configuration
nav_order: 3
---

## Caching in NaSe

A CLI flag has been added to optionally enable caching for the nameservice.
Pass `--nase-cached` to your `fred` instance to activate caching.
This improves performance for requests to `fred` but may lead to data inconsistency if configuration changes often.
By default, it is turned off.
