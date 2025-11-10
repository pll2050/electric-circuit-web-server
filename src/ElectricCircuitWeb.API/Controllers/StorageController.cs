using Microsoft.AspNetCore.Mvc;
using ElectricCircuitWeb.API.Services;

namespace ElectricCircuitWeb.API.Controllers;

/// <summary>
/// 스토리지 컨트롤러
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class StorageController : ControllerBase
{
    private readonly IStorageService _storageService;
    private readonly ILogger<StorageController> _logger;

    public StorageController(IStorageService storageService, ILogger<StorageController> logger)
    {
        _storageService = storageService;
        _logger = logger;
    }

    // ==================== 파일 업로드 ====================

    /// <summary>
    /// 파일을 서버에 업로드합니다.
    /// POST /api/storage/upload
    /// </summary>
    [HttpPost("upload")]
    [Consumes("multipart/form-data")]
    public async Task<IActionResult> UploadFile(IFormFile file, string? folder = null)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (file == null || file.Length == 0)
            {
                return BadRequest(new
                {
                    success = false,
                    error = "No file uploaded"
                });
            }

            // 파일 크기 체크 (32MB)
            if (file.Length > 32 * 1024 * 1024)
            {
                return BadRequest(new
                {
                    success = false,
                    error = "File size exceeds 32MB limit"
                });
            }

            using var stream = file.OpenReadStream();
            var result = await _storageService.UploadFileAsync(stream, file.FileName, folder);

            return Ok(new
            {
                success = true,
                message = "File uploaded successfully",
                download_url = result.DownloadUrl,
                file_path = result.FilePath,
                file_name = result.FileName,
                size = result.Size
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error uploading file");
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 파일 URL 조회 ====================

    /// <summary>
    /// 저장된 파일의 다운로드 URL을 조회합니다.
    /// GET /api/storage/url?filePath={file_path}
    /// </summary>
    [HttpGet("url")]
    public async Task<IActionResult> GetFileUrl([FromQuery] string filePath)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (string.IsNullOrEmpty(filePath))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "File path is required"
                });
            }

            var url = await _storageService.GetFileUrlAsync(filePath);

            return Ok(new
            {
                success = true,
                message = "File URL retrieved successfully",
                download_url = url
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting file URL: {FilePath}", filePath);
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 파일 삭제 ====================

    /// <summary>
    /// 저장된 파일을 삭제합니다.
    /// DELETE /api/storage/delete?filePath={file_path}
    /// </summary>
    [HttpDelete("delete")]
    public async Task<IActionResult> DeleteFile([FromQuery] string filePath)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (string.IsNullOrEmpty(filePath))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "File path is required"
                });
            }

            var deleted = await _storageService.DeleteFileAsync(filePath);

            if (!deleted)
            {
                return NotFound(new
                {
                    success = false,
                    error = "File not found"
                });
            }

            return Ok(new
            {
                success = true,
                message = "File deleted successfully"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error deleting file: {FilePath}", filePath);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 파일 목록 조회 ====================

    /// <summary>
    /// 사용자의 파일 목록을 조회합니다.
    /// GET /api/storage/list?folder={folder_name}
    /// </summary>
    [HttpGet("list")]
    public async Task<IActionResult> ListFiles([FromQuery] string? folder = null)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            var files = await _storageService.ListFilesAsync(folder);

            return Ok(new
            {
                success = true,
                message = "Files listed successfully",
                files = files.Select(f => new
                {
                    name = f.Name,
                    path = f.Path,
                    size = f.Size,
                    type = f.Type,
                    url = f.Url,
                    created_at = f.CreatedAt
                })
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error listing files");
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 회로 이미지 업로드 ====================

    /// <summary>
    /// 회로 이미지를 업로드합니다.
    /// POST /api/storage/upload-circuit-image
    /// </summary>
    [HttpPost("upload-circuit-image")]
    [Consumes("multipart/form-data")]
    public async Task<IActionResult> UploadCircuitImage(IFormFile image, string circuitId)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (image == null || image.Length == 0)
            {
                return BadRequest(new
                {
                    success = false,
                    error = "No image uploaded"
                });
            }

            if (string.IsNullOrEmpty(circuitId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Circuit ID is required"
                });
            }

            // 이미지 파일 확인
            var allowedTypes = new[] { "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp" };
            if (!allowedTypes.Contains(image.ContentType.ToLower()))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Only image files are allowed (JPEG, PNG, GIF, WebP)"
                });
            }

            // 파일 크기 체크 (32MB)
            if (image.Length > 32 * 1024 * 1024)
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Image size exceeds 32MB limit"
                });
            }

            using var stream = image.OpenReadStream();
            var result = await _storageService.UploadCircuitImageAsync(stream, image.FileName, circuitId);

            return Ok(new
            {
                success = true,
                message = "Circuit image uploaded successfully",
                download_url = result.DownloadUrl,
                file_path = result.FilePath
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error uploading circuit image");
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }
}
