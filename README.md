# SL2-hash
### Хэширование данный при помощи SL2

Алгоритм основан на построении группы *SL2(F_p)*, элементами которой являются матрицы 2 на 2. 

https://link.springer.com/content/pdf/10.1007/3-540-48658-5_5.pdf

## Install

     go get github.com/qwertyqq2/sl2-task


## Usage

     // создание генетора группы SL2
     gen := sl2.Generate(
		sl2.SetOrderField128(),
		sl2.SetDefaultElement(),
		sl2.SetSha256(),
	)

     // цепочка данных
     
     data := "my name is shao khan"
	
     strs := strings.Split(data, " ") 

     // хэш
     hash, err := gen.Snapshot(strs...) // хэш
	


