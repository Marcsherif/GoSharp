@fmt = private constant [4 x i8] c"%d\0A\00"

declare i32 @printf(ptr, ...)

define i32 @main() {
entry:
  %x.addr = alloca i32
  store i32 0, ptr %x.addr
  %x.value.0 = load i32, ptr %x.addr
  call i32 (ptr, ...) @printf(ptr @fmt, i32 %x.value.0)
  ret i32 0
}
