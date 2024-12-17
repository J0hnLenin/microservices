using System;
using System.IO;
using System.Net.Http;
using System.Threading.Tasks;

class Program
{
    static async Task Main(string[] args)
    {
        Task task1 = RunAsyncService(1);
        Task task2 = RunAsyncService(2);
        await Task.WhenAll(task1, task2);
    }

    static async Task RunAsyncService(int id)
    {
        string getImageUrl = "http://go-service:8081/generate-image";
        string postCompressURL = "http://nodejs-service:8082/compress";
        string postPredictURL = "http://python-service:8083/predict";
        string size = "128";

        while (true)
        {
            using (HttpClient client = new HttpClient())
            {
                byte[] imageBytes;
                byte[] imageBytesCompressed;
                client.DefaultRequestHeaders.Add("size", size);
                Console.WriteLine($"{DateTime.Now} INFO: task {id}; send request to {getImageUrl}");
                try
                {
                    HttpResponseMessage response = await client.GetAsync(getImageUrl);
                    response.EnsureSuccessStatusCode();
                    imageBytes = await response.Content.ReadAsByteArrayAsync();
                    Console.WriteLine($"{DateTime.Now} INFO: task {id}; code {response.StatusCode}; image received; size {imageBytes.Length} bytes");
                }
                catch (Exception e)
                {
                    Console.WriteLine($"{DateTime.Now} ERROR: task {id}; {e.Message}");
                    await Task.Delay(5000);
                    continue;
                }
                Console.WriteLine($"{DateTime.Now} INFO: task {id}; send request to {postCompressURL}");
                try
                {
                    var content = new MultipartFormDataContent();
                    var byteArrayContent = new ByteArrayContent(imageBytes);
                    content.Add(byteArrayContent, "image", "image.png");

                    HttpResponseMessage response = await client.PostAsync(postCompressURL, content);
                    response.EnsureSuccessStatusCode();

                    imageBytesCompressed = await response.Content.ReadAsByteArrayAsync();
                    Console.WriteLine($"{DateTime.Now} INFO: task {id}; code {response.StatusCode}; image received; size {imageBytesCompressed.Length} bytes");
                }
                catch (Exception e)
                {
                    Console.WriteLine($"{DateTime.Now} ERROR: task {id}; {e.Message}");
                    await Task.Delay(5000);
                    continue;
                }
                Console.WriteLine($"{DateTime.Now} INFO: task {id}; send request to {postPredictURL}");
                try
                {
                    var content = new MultipartFormDataContent();
                    var byteArrayContent = new ByteArrayContent(imageBytesCompressed);
                    content.Add(byteArrayContent, "image", "image.png");

                    HttpResponseMessage response = await client.PostAsync(postPredictURL, content);
                    response.EnsureSuccessStatusCode();
                    string responseBody = await response.Content.ReadAsStringAsync();
                    Console.WriteLine($"{DateTime.Now} INFO: task {id}; code {response.StatusCode}; answer {responseBody}");
                }
                catch (Exception e)
                {
                    Console.WriteLine($"{DateTime.Now} ERROR: task {id}; {e.Message}");
                    await Task.Delay(5000);
                    continue;
                }
            }
            await Task.Delay(5000);
        }
    }
}