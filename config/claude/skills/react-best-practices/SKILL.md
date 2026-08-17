---
name: react-best-practices
description: Provides React patterns for hooks, effects, refs, and component design. Covers escape hatches, anti-patterns, and correct effect usage. Must use when reading or writing React components (.tsx, .jsx files with React imports).
---

# React Best Practices

## Core Principle: Effects Are Escape Hatches

Effects let you "step outside" React to synchronize with external systems. **Most component logic should NOT use Effects.** Before writing an Effect, ask: "Is there a way to do this without an Effect?"

## When to Use Effects

Effects are for synchronizing with **external systems**:
- Subscribing to browser APIs (WebSocket, IntersectionObserver, resize)
- Connecting to third-party libraries not written in React
- Setting up/cleaning up event listeners on window/document
- Fetching data on mount (though prefer React Query or framework data fetching)
- Controlling non-React DOM elements (video players, maps, modals)

## When NOT to Use Effects

### Derived State (Calculate During Render)

```tsx
// BAD: Effect for derived state
const [firstName, setFirstName] = useState('Taylor');
const [lastName, setLastName] = useState('Swift');
const [fullName, setFullName] = useState('');
useEffect(() => {
  setFullName(firstName + ' ' + lastName);
}, [firstName, lastName]);

// GOOD: Calculate during render
const [firstName, setFirstName] = useState('Taylor');
const [lastName, setLastName] = useState('Swift');
const fullName = firstName + ' ' + lastName;
```

### Expensive Calculations (Use useMemo)

```tsx
// BAD: Effect for caching
const [visibleTodos, setVisibleTodos] = useState([]);
useEffect(() => {
  setVisibleTodos(getFilteredTodos(todos, filter));
}, [todos, filter]);

// GOOD: useMemo for expensive calculations
const visibleTodos = useMemo(
  () => getFilteredTodos(todos, filter),
  [todos, filter]
);
```

### Resetting State on Prop Change (Use key)

```tsx
// BAD: Effect to reset state
function ProfilePage({ userId }) {
  const [comment, setComment] = useState('');
  useEffect(() => {
    setComment('');
  }, [userId]);
  // ...
}

// GOOD: Use key to reset component state
function ProfilePage({ userId }) {
  return <Profile userId={userId} key={userId} />;
}

function Profile({ userId }) {
  const [comment, setComment] = useState(''); // Resets automatically
  // ...
}
```

### User Event Handling (Use Event Handlers)

```tsx
// BAD: Event-specific logic in Effect
function ProductPage({ product, addToCart }) {
  useEffect(() => {
    if (product.isInCart) {
      showNotification(`Added ${product.name} to cart`);
    }
  }, [product]);
  // ...
}

// GOOD: Logic in event handler
function ProductPage({ product, addToCart }) {
  function buyProduct() {
    addToCart(product);
    showNotification(`Added ${product.name} to cart`);
  }
  // ...
}
```

### Notifying Parent of State Changes

```tsx
// BAD: Effect to notify parent
function Toggle({ onChange }) {
  const [isOn, setIsOn] = useState(false);
  useEffect(() => {
    onChange(isOn);
  }, [isOn, onChange]);
  // ...
}

// GOOD: Update both in event handler
function Toggle({ onChange }) {
  const [isOn, setIsOn] = useState(false);
  function updateToggle(nextIsOn) {
    setIsOn(nextIsOn);
    onChange(nextIsOn);
  }
  // ...
}

// BEST: Fully controlled component
function Toggle({ isOn, onChange }) {
  function handleClick() {
    onChange(!isOn);
  }
  // ...
}
```

### Chains of Effects

```tsx
// BAD: Effect chain
useEffect(() => {
  if (card !== null && card.gold) {
    setGoldCardCount(c => c + 1);
  }
}, [card]);

useEffect(() => {
  if (goldCardCount > 3) {
    setRound(r => r + 1);
    setGoldCardCount(0);
  }
}, [goldCardCount]);

// GOOD: Calculate derived state, update in event handler
const isGameOver = round > 5;

function handlePlaceCard(nextCard) {
  setCard(nextCard);
  if (nextCard.gold) {
    if (goldCardCount < 3) {
      setGoldCardCount(goldCardCount + 1);
    } else {
      setGoldCardCount(0);
      setRound(round + 1);
    }
  }
}
```

## Summary: Decision Tree

1. **Need to respond to user interaction?** Use event handler
2. **Need computed value from props/state?** Calculate during render
3. **Need cached expensive calculation?** Use useMemo
4. **Need to reset state on prop change?** Use key prop
5. **Need to synchronize with external system?** Use Effect with cleanup
6. **Need non-reactive code in Effect?** Use useEffectEvent
7. **Need mutable value that doesn't trigger render?** Use ref


## Detailed references

Load these as needed:

- **Effect dependencies**: `references/effect-dependencies.md`. Dependency arrays, stale closures, removing unnecessary dependencies.
- **Effect cleanup**: `references/effect-cleanup.md`. Cleanup functions, race conditions, subscriptions.
- **Refs**: `references/refs.md`. When to use a ref instead of state, DOM refs, ref callbacks.
- **Custom hooks**: `references/custom-hooks.md`. Extracting logic, naming, what belongs in a hook.
- **Component patterns**: `references/component-patterns.md`. Composition, keys, lifting state, list rendering.
