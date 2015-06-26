# 法務部行政執行署動產拍賣爬蟲程式
*法務部行政執行署動產拍賣爬蟲程式*

## 用法
```
import (
	"fmt"
	"mp"
)

func main() {
	s := mp.NewService()
	s.Get()
	fmt.Println(len(s.Items))
	// 顯示所有拍賣數量
}
```

## 許可證
mp 使用 CC0 許可證，詳參 LICENSE 檔案。