# album-maker

<p align="left">
  <a href="https://github.com/5hyn3/album-maker/actions?query=workflow%3AGo"><img alt="GitHub Actions status" src="https://github.com/5hyn3/album-maker/workflows/Go/badge.svg"></a>
</p>

album maker 

```
target-dir/
├── a // Last modified 2019/11/11
├── b // Last modified 2019/11/11
├── c // Last modified 2019/11/11
└── d // Last modified 2019/11/12
```

```
album-maker --target-dir=target-dir
```


```
target-dir/
└── 2019
    └── 11
        ├── 11
        │   ├── a
        │   ├── b
        │   └── c
        └── 12
            └── d
```
