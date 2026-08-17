## Component Patterns

### Controlled vs Uncontrolled

```tsx
// Uncontrolled: component owns state
function SearchInput() {
  const [query, setQuery] = useState('');
  return <input value={query} onChange={e => setQuery(e.target.value)} />;
}

// Controlled: parent owns state
function SearchInput({ query, onQueryChange }) {
  return <input value={query} onChange={e => onQueryChange(e.target.value)} />;
}
```

### Prefer Composition Over Prop Drilling

```tsx
// BAD: Prop drilling
<App user={user}>
  <Layout user={user}>
    <Header user={user}>
      <Avatar user={user} />
    </Header>
  </Layout>
</App>

// GOOD: Composition with children
<App>
  <Layout>
    <Header avatar={<Avatar user={user} />} />
  </Layout>
</App>

// GOOD: Context for truly global state
<UserContext.Provider value={user}>
  <App />
</UserContext.Provider>
```

### flushSync for Synchronous DOM Updates

```tsx
// When you need to read DOM immediately after state update
import { flushSync } from 'react-dom';

function handleAdd() {
  flushSync(() => {
    setTodos([...todos, newTodo]);
  });
  // DOM is now updated, safe to read
  listRef.current.lastChild.scrollIntoView();
}
```

