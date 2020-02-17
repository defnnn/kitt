from pydantic import BaseModel


class Example(BaseModel):
    foo: str
    bar: str


def sm(a: int, b: int) -> int:
    return a + b
