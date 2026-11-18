@fmt = private constant [4 x i8] c"%d\0A\00"

declare i32 @printf(ptr, ...)

define i32 @main() {
entry:
  %x.addr = alloca i32
  store i32 0, ptr %x.addr

  %result = load i32, ptr %x.addr
  call i32 (ptr, ...) @printf(ptr @fmt, i32 %result)
  ret i32 0
}
