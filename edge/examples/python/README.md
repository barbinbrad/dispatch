# Running the Example

First start the proxy server. Then run the example with the following command:

```sh
python3 device.py
```

The example process implements two methods. `example` and `exampleWithParams`.

```python
@dispatcher.add_method
def example():
  return {"success": 1}

@dispatcher.add_method
def exampleWithParams(sleep):
  time.sleep(sleep)
  return {
    "slept": sleep
  }


```