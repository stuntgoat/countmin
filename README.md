Countmin sketch as translated from [http://www.cs.rutgers.edu/~muthu/massdal-code-index.html](http://www.cs.rutgers.edu/~muthu/massdal-code-index.html). Usage:

    func main () {
        rand.Seed( time.Now().UTC().UnixNano())
        // 1200 columns
        width := 5000
        // 10 rows
        depth := 10
        // Specify width and depth
        cm := MakeCM(5000, 10)

        // For each i, increment the sketch at i to 1
        var i int64
        for i = 0; i < 2500; i++ {
            cm.Update(i, 1)
        }

        // Check that each j < 2500 == 1
        var val int64
        var j int64
        for j = 0; j < 5000; j++ {
            val = cm.PointEst(j)
            fmt.Println("point estimate for", j, val)
    }
