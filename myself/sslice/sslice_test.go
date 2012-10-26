package sslice_test

import (
    "fmt"
    "github.com/tw4452852/Programming_in_go-exercise/myself/sslice"
    "sync"
    "testing"
)

func TestSafeSlice(t *testing.T) {
    store :=sslice.New()
    fmt.Printf("Initially has %d items\n", store.Len())

    deleted := []int{0, 2, 3, 5, 7, 20, 399, 25, 30, 1000, 91, 97, 98, 99}

    var waiter sync.WaitGroup

    waiter.Add(1)
    go func() { // Concurrent Inserter
        for i := 0; i < 100; i++ {
            store.Append(fmt.Sprintf("%04X", i))
            if i > 0 && i%15 == 0 {
                fmt.Printf("Inserted %d items\n", store.Len())
            }
        }
        fmt.Printf("Inserted %d items\n", store.Len())
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Deleter
        for _, i := range deleted {
            before := store.Len()
            store.Delete(i)
            fmt.Printf("Deleted m[%d] before=%d after=%d\n",
                i, before, store.Len())
        }
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Finder
        for _, i := range deleted {
            for _, j := range []int{i, i + 1} {
                item := store.At(j)
                if item != nil {
                    fmt.Printf("Found m[%d] == %s\n", j, item)
                } else {
                    fmt.Printf("Not found m[%d]\n", j)
                }
            }
        }
        waiter.Done()
    }()

    waiter.Wait()

    fmt.Printf("Finished with %d items\n", store.Len())
    updater := func(value interface{}) interface{} {
        return value.(string) + ":updated"
    }
    for i := 0; i < store.Len() && i < 5; i++ {
        fmt.Printf("m[%d] == %s -> ", i, store.At(i))
        store.Update(i, updater)
        fmt.Printf("%s\n", store.At(i))
    }
    list := store.Close()
    fmt.Println("Closed")
    fmt.Printf("len == %d\n", len(list))
    fmt.Println()
}
