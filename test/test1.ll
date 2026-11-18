@fmt = private constant [4 x i8] c"%d\0A\00"

declare i32 @printf(ptr, ...)

define i32 @main() {
entry:
  %x.addr = alloca i32
  %i.addr = alloca i32
  store i32 0, ptr %x.addr
  store i32 1, ptr %i.addr
  br label %while.cond

while.cond:
  %i.current = load i32, ptr %i.addr
  %keep.going = icmp slt i32 %i.current, 10
  br i1 %keep.going, label %while.body, label %while.end

while.body:
  %x.current = load i32, ptr %x.addr
  %square = mul i32 %i.current, %i.current
  %sum = add i32 %x.current, %square
  store i32 %sum, ptr %x.addr
  %next.i = add i32 %i.current, 1
  store i32 %next.i, ptr %i.addr
  br label %while.cond

while.end:
  %result = load i32, ptr %x.addr
  call i32 (ptr, ...) @printf(ptr @fmt, i32 %result)
  ret i32 0
}
