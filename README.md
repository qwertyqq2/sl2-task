# SL2-hash
### Хэширование данный при помощи SL2

Алгоритм основан на построении группы *SL2(F_p)*, элементами которой являются матрицы 2 на 2. 

## Install

     go get github.com/qwertyqq2/sl2-task


## Usage

     gen := sl2.Generate(
		sl2.SetOrderField128(),
		sl2.SetDefaultElement(),
		sl2.SetSha256(),
	) // создание генетора группы *SL2*

     data := "my name is shao khan"
	strs := strings.Split(data, " ") //цепочка данных

     hash, err := gen.Snapshot(strs...) // хэш
	


