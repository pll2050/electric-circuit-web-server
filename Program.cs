using Microsoft.EntityFrameworkCore;
using Microsoft.AspNetCore.Http;
using FirebaseAdmin;
using Google.Apis.Auth.OAuth2;
using ElectricCircuitWeb.API.Data;
using ElectricCircuitWeb.API.Services;
using ElectricCircuitWeb.API.Repositories;
using ElectricCircuitWeb.API.Config;
using ElectricCircuitWeb.API.Middleware;

var builder = WebApplication.CreateBuilder(args);

// Configuration
var configuration = builder.Configuration;

// Add PostgreSQL
var connectionString = configuration.GetConnectionString("DefaultConnection")
    ?? "Host=localhost;Port=5432;Database=electric_circuit_db;Username=postgres;Password=q1w2e3r4";

builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseNpgsql(connectionString));

// Add Firebase
var firebaseConfig = configuration.GetSection("Firebase").Get<FirebaseConfig>();
if (!string.IsNullOrEmpty(firebaseConfig?.ProjectId))
{
    var credential = GoogleCredential.FromFile(firebaseConfig.ServiceAccountKeyPath);
    FirebaseApp.Create(new AppOptions
    {
        Credential = credential,
        ProjectId = firebaseConfig.ProjectId
    });
}

// Add CORS
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyMethod()
              .AllowAnyHeader();
    });
});

// Add services
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(options =>
{
    // Support for file uploads in Swagger
    options.MapType<IFormFile>(() => new Microsoft.OpenApi.Models.OpenApiSchema
    {
        Type = "string",
        Format = "binary"
    });
});

// Register repositories
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddScoped<IProjectRepository, ProjectRepository>();
builder.Services.AddScoped<ICircuitRepository, CircuitRepository>();

// Register services
builder.Services.AddScoped<IAuthService, AuthService>();
builder.Services.AddScoped<IProjectService, ProjectService>();
builder.Services.AddScoped<ICircuitService, CircuitService>();
builder.Services.AddScoped<IStorageService, StorageService>();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

// Middleware
app.UseCors();
app.UseHttpsRedirection();

// Firebase Auth Middleware (optional, 주석 처리됨)
// app.UseFirebaseAuth();

app.UseAuthorization();
app.MapControllers();

// Database migration (개발 환경에서만)
if (app.Environment.IsDevelopment())
{
    using var scope = app.Services.CreateScope();
    var dbContext = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
    try
    {
        await dbContext.Database.MigrateAsync();
        app.Logger.LogInformation("Database migrated successfully");
    }
    catch (Exception ex)
    {
        app.Logger.LogError(ex, "Failed to migrate database");
    }
}

app.Run();
