# go_api_cache_prac2

参考：https://golang.hateblo.jp/entry/golang-http-cache

API結果をメモリにキャッシュ、
2回目以降のアクセスではetagとlast-modifiedを使って変更の有無の確認をしたうえで、キャッシュを使うか処理を分けている
