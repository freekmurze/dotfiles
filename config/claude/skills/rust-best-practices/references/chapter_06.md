# Chapter 6 - Generics, Dynamic Dispatch and Static Dispatch

> Static where you can, dynamic where you must

Rust allows you to handle polymorphic code in two ways:
* **Generics / Static Dispatch**: compile-time, monomorphized per use.
* **Trait Objects / Dynamic Dispatch**: runtime vtable, single implementation.

Understanding the trade-offs lets you write faster, smaller and more flexible code.

## 6.1 [Generics](https://doc.rust-lang.org/book/ch10-00-generics.html)

Every programming language has tools for effectively handling the duplication of concepts. In Rust, one such tool is generics: abstract stand-ins for concrete types or other properties. We can express the behavior of generics or how they relate to other generics without knowing what will be in their place when compiling and running the code. 

We use generics to create definitions for items like function signatures or structs, which we can then use with many different concrete data types. Let's first look at how to define functions, structs, enums, and methods using generics. Generics can also be used to implement Type State Pattern and constrain a struct functionality to certain expected types, more on type state on [chapter 7](./chapter_07.md).

[Generics by Examples](https://doc.rust-lang.org/rust-by-example/generics.html).

### Generics Performance

You might be wondering whether there is a runtime cost when using generic type parameters. The good news is that using generic types won't make your program run any slower than it would with concrete types. Rust accomplishes this by performing monomorphization of the code using generics at compile time. Monomorphization is the process of turning generic code into specific code by filling in the concrete types that are used when compiled. The compiler checks for all occurrences of the generic parameter and generates code for the concrete types the generic code is called with.

## 6.2 Static Dispatch: `impl Trait` or `<T: Trait>`

A static dispatch is basically a constrained version of generics, a trait bounded generic, at compile-time it is able to check if your generic satisfies the declared traits.

### ✅  Best when:
* You want **zero runtime cost**, by paying the compile time cost.
* You need **tight loops or performance**.
* Your types are **known at compile time**.
* You are working with **single-use implementations** (monomorphized).

### 🏎️ Example: High-performance function with generic
```rust
fn specialized_sum<U: Sum + RandomMapping, T: Iterator<Item = U>>(iter: T) -> U {
    iter.map(|x| x.random_mapping()).sum()
}

// or, equivalent, more modern
fn specialized_sum<U: Sum + RandomMapping>(iter: impl Iterator<Item = U>) -> U {
    iter.map(|x| x.random_mapping()).sum()
}
```

This is compiled into **specialized machine code** for each usage, fast and inlined.

## 6.3 Dynamic Dispatch: `dyn Trait`

Usually dynamic dispatch is used with some kind of pointer or a reference, like `Box<dyn Trait>`, `Arc<dyn Trait>` or `&dyn trait`.

### ✅  Best when:
* You absolutely need runtime polymorphism.
* You need to **store different implementations** in one collection.
* You want to **abstract internals behind a stable interface**.
* You are writing a **plugin-style architecture**.

> ❗ Closer to what you would get in an object oriented language and can have some heavy costs associated to it. Can avoid generics entirely and let you mix types that implement the same traits.

### 🚚 Example: Heterogeneous collection

```rust
trait Animal {
    fn greet(&self) -> String;
}

struct Dog;
impl Animal for Dog {
    fn greet(&self) -> String {
        "woof".to_string()
    }
}

struct Cat;
impl Animal for Cat {
    fn greet(&self) -> String {
        "meow".to_string()
    }
}

fn all_animals_greeting(animals: Vec<Box<dyn Animal>>) {
    for animal in animals {
        println!("{}", animal.greet())
    }
}
```

## 6.4 Trade-off summary

|                   	| Static Dispatch (impl Trait) 	|    Dynamic Dispatch (dyn Trait)   	|
|-------------------	|------------------------------	|---------------------------------- 	|
| Performance       	| ✅ Faster, inlined            	| ❌ Slower: vtable indirection         |
| Compile time      	| ❌ Slower: monomorphization   	| ✅ Faster: shared code                |
| Binary size       	| ❌ Larger: per-type codegen   	| ✅ Smaller                            |
| Flexibility       	| ❌ Rigid, one type at a time  	| ✅ Can mix types in collections       |
| Use in trait fn() 	| ❌ Traits must be object-safe 	| ✅ Works with trait objects           |
| Errors            	| ✅ Clearer                    	| ❌ Erased types can confuse errors    |

* Prefer generics/static dispatch when you control the call site and want performance.
* Use dynamic dispatch when you need abstraction, plugins or mixed types. 🚨 Runtime cost.
* If you are not sure, start with generics, trait bound them - then use `Box<dyn Trait>` when flexibility outweighs speed.

> Favor static dispatch until your trait needs to live behind a pointer.

## 6.5 Best Practices for Dynamic Dispatch

Dynamic dispatch `Ptr<dyn Trait>` is a powerful tool, but it also has significant performance trade-offs. You should only reach for it when **type erasure or runtime polymorphism** are essential. It is important to know when you need Trait Objects:

### ✅ Use Dynamic Dispatch When:

* You need heterogeneous types in a collection:
```rust
fn all_animals_greeting(animals: Vec<Box<dyn Animal>>) {
    for animal in animals {
        println!("{}", animal.greet())
    }
}
```

* You want runtime plugins or hot-swappable components.
* You want to abstract internals from the caller (library design).


### ❌ Avoid Dynamic Dispatch When:

* You control the concrete types.
* You are writing code in performance critical paths.
* You can express the same logic in other ways while keeping simplicity, e.g. generics.

## 6.6 🚨 Trait Objects Ergonomics

* Prefer `&dyn Trait` over `Box<dyn Trait>` when you don't need ownership.
* Use `Arc<dyn Trait + Send + Sync>` for shared access across threads.
* Don't use `dyn Trait` if the trait has methods that return `Self`.
* **Avoid boxing too early**. Don't box inside structs unless you are sure it'll be beneficial or is required (recursive).
```rust
// ✅ Use generics when possible
struct Renderer<B: RenderBackend> {
    backend: B
}

// ❌ Premature Boxing
struct Renderer {
    backend: Box<dyn RenderBackend> // Boxing too early
}
```
* If you must expose a `dyn trait` in a public API, `Box` at the boundary, not internally.
* **Object Safety**: You can only create `dyn Traits` from object-safe traits:
    * It has **no generic methods**.
    * It doesn't require `Self: Sized`.
    * All method signatures use `&self`, `&mut self` or `self`.
    ```rust
    // ✅ Object Safe
    trait Runnable {
        fn run(&self);
    }

    // ❌ Not Object Safe
    trait Factory {
        fn create<T>() -> T; // generic methods are not allowed
    }
    ```