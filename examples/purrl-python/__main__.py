import pulumiverse_purrl as purrl
import pulumi

purrl_command = purrl.Purrl(
    "purrl-python",
    name="purrl-python",
    method="GET",
    headers={"test": "test"},
    url="https://httpbin.org/get",
    expected_response_codes=["200"],
    delete_method="DELETE",
    delete_url="https://httpbin.org/delete",
    expected_delete_response_codes=["200"],
)

pulumi.export("response", purrl_command.response)
