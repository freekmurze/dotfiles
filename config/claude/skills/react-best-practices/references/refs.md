## Refs

### Use Refs for Values That Don't Affect Rendering

```tsx
// GOOD: Ref for timeout ID (doesn't affect UI)
const timeoutRef = useRef(null);

function handleClick() {
  clearTimeout(timeoutRef.current);
  timeoutRef.current = setTimeout(() => {
    // ...
  }, 1000);
}

// BAD: Using ref for displayed value
const countRef = useRef(0);
countRef.current++; // UI won't update!
```

### Never Read/Write ref.current During Render

```tsx
// BAD: Reading ref during render
function MyComponent() {
  const ref = useRef(0);
  ref.current++; // Mutating during render!
  return <div>{ref.current}</div>; // Reading during render!
}

// GOOD: Read/write refs in event handlers and effects
function MyComponent() {
  const ref = useRef(0);

  function handleClick() {
    ref.current++; // OK in event handler
  }

  useEffect(() => {
    ref.current = someValue; // OK in effect
  }, [someValue]);
}
```

### Ref Callbacks for Dynamic Lists

```tsx
// BAD: Can't call useRef in a loop
{items.map((item) => {
  const ref = useRef(null); // Rule violation!
  return <li ref={ref} />;
})}

// GOOD: Ref callback with Map
const itemsRef = useRef(new Map());

{items.map((item) => (
  <li
    key={item.id}
    ref={(node) => {
      if (node) {
        itemsRef.current.set(item.id, node);
      } else {
        itemsRef.current.delete(item.id);
      }
    }}
  />
))}
```

### useImperativeHandle for Controlled Exposure

```tsx
// Limit what parent can access
function MyInput({ ref }) {
  const realInputRef = useRef(null);

  useImperativeHandle(ref, () => ({
    focus() {
      realInputRef.current.focus();
    },
    // Parent can ONLY call focus(), not access full DOM node
  }));

  return <input ref={realInputRef} />;
}
```

