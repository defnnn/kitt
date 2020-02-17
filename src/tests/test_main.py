from main import Example, sm


def test_ex(ex: Example) -> None:
    foo = "FOO"
    bar = "BAR"
    example = Example(foo=foo, bar=bar)
    assert example.foo == foo
    assert example.bar == bar


def test_sm() -> None:
    a = 1
    b = 2
    assert sm(a, b) == a + b
