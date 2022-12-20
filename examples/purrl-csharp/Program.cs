using System.Collections.Generic;
using Pulumi;
using Pulumiverse.Purrl;

return await Deployment.RunAsync(() =>
{
   var purrl =new Purrl("purrl", new PurrlArgs
   {
      Name = "httpbin",
      Url = "https://httpbin.org/get",
      ResponseCodes = new List<string> { "200" },
      Method = "GET",
      Headers = new Dictionary<string, string> { { "test", "test" } },
      DeleteMethod = "DELETE",
      DeleteUrl = "https://httpbin.org/delete",
      DeleteResponseCodes = new List<string> { "200" },
   });

   // Export outputs here
   return new Dictionary<string, object?>
   {
      ["response"] =purrl.Response
   };
});
