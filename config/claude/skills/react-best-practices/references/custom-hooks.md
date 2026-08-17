## Custom Hooks

### Hooks Share Logic, Not State

```tsx
// Each call gets independent state
function StatusBar() {
  const isOnline = useOnlineStatus(); // Own state
}

function SaveButton() {
  const isOnline = useOnlineStatus(); // Separate state instance
}
```

### Name Hooks useXxx Only If They Use Hooks

```tsx
// BAD: useXxx but doesn't use hooks
function useSorted(items) {
  return items.slice().sort();
}

// GOOD: Regular function
function getSorted(items) {
  return items.slice().sort();
}

// GOOD: Uses hooks, so prefix with use
function useAuth() {
  return useContext(AuthContext);
}
```

### Avoid "Lifecycle" Hooks

```tsx
// BAD: Custom lifecycle hooks
function useMount(fn) {
  useEffect(() => {
    fn();
  }, []); // Missing dependency, linter can't catch it
}

// GOOD: Use useEffect directly
useEffect(() => {
  doSomething();
}, [doSomething]);
```

### Keep Custom Hooks Focused

```tsx
// GOOD: Focused, concrete use cases
useChatRoom({ serverUrl, roomId });
useOnlineStatus();
useFormInput(initialValue);

// BAD: Generic, abstract hooks
useMount(fn);
useEffectOnce(fn);
useUpdateEffect(fn);
```

