namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 스토리지 서비스 인터페이스
/// </summary>
public interface IStorageService
{
    // 파일 업로드
    Task<StorageUploadResult> UploadFileAsync(Stream fileStream, string fileName, string? folder = null);

    // 파일 URL 조회
    Task<string> GetFileUrlAsync(string filePath);

    // 파일 삭제
    Task<bool> DeleteFileAsync(string filePath);

    // 파일 목록 조회
    Task<IEnumerable<StorageFileInfo>> ListFilesAsync(string? folder = null);

    // 회로 이미지 업로드
    Task<StorageUploadResult> UploadCircuitImageAsync(Stream imageStream, string fileName, string circuitId);
}

/// <summary>
/// 스토리지 업로드 결과
/// </summary>
public class StorageUploadResult
{
    public string DownloadUrl { get; set; } = string.Empty;
    public string FilePath { get; set; } = string.Empty;
    public string FileName { get; set; } = string.Empty;
    public long Size { get; set; }
}

/// <summary>
/// 스토리지 파일 정보
/// </summary>
public class StorageFileInfo
{
    public string Name { get; set; } = string.Empty;
    public string Path { get; set; } = string.Empty;
    public long Size { get; set; }
    public string Type { get; set; } = string.Empty;
    public string Url { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; }
}
