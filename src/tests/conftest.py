import pytest

from main import Example


@pytest.fixture(scope="module")
def ex():
    ex_fixture = Example(foo="foo...", bar="bar...")
    yield ex_fixture
