# csv
导入csv数据结构

## rock.csv
- v = vela.csv(file)
- file 文件路径

#### 内部方法
- [v.pipe(object)]()
```lua
    local c = rock.csv("a.csv")
    c.pipe(function(row)
        print(row[1])
        print(row[2])
    end)
```