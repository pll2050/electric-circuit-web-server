namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 스토리지 서비스 구현 (Firebase Storage)
/// </summary>
public class StorageService : IStorageService
{
    private readonly ILogger<StorageService> _logger;

    public StorageService(ILogger<StorageService> logger)
    {
        _logger = logger;
    }

    public async Task<StorageUploadResult> UploadFileAsync(Stream fileStream, string fileName, string? folder = null)
    {
        // TODO: Firebase Storage 구현
        _logger.LogWarning("Firebase Storage not yet implemented");

        await Task.CompletedTask;

        return new StorageUploadResult
        {
            DownloadUrl = $"https://storage.example.com/files/{fileName}",
            FilePath = folder != null ? $"{folder}/{fileName}" : fileName,
            FileName = fileName,
            Size = fileStream.Length
        };
    }

    public async Task<string> GetFileUrlAsync(string filePath)
    {
        // TODO: Firebase Storage 구현
        _logger.LogWarning("Firebase Storage not yet implemented");

        await Task.CompletedTask;

        return $"https://storage.example.com/files/{filePath}";
    }

    public async Task<bool> DeleteFileAsync(string filePath)
    {
        // TODO: Firebase Storage 구현
        _logger.LogWarning("Firebase Storage not yet implemented");

        await Task.CompletedTask;

        return true;
    }

    public async Task<IEnumerable<StorageFileInfo>> ListFilesAsync(string? folder = null)
    {
        // TODO: Firebase Storage 구현
        _logger.LogWarning("Firebase Storage not yet implemented");

        await Task.CompletedTask;

        return new List<StorageFileInfo>();
    }

    public async Task<StorageUploadResult> UploadCircuitImageAsync(Stream imageStream, string fileName, string circuitId)
    {
        // TODO: Firebase Storage 구현
        _logger.LogWarning("Firebase Storage not yet implemented");

        var folder = $"circuits/{circuitId}";
        return await UploadFileAsync(imageStream, fileName, folder);
    }
}
