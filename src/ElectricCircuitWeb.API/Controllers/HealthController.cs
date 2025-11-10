using Microsoft.AspNetCore.Mvc;

namespace ElectricCircuitWeb.API.Controllers;

/// <summary>
/// 헬스 체크 컨트롤러
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class HealthController : ControllerBase
{
    /// <summary>
    /// 서버 상태를 확인합니다.
    /// </summary>
    [HttpGet]
    public IActionResult Get()
    {
        return Ok(new
        {
            status = "healthy",
            timestamp = DateTime.UtcNow
        });
    }
}
