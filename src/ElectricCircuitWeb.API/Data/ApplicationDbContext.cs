using Microsoft.EntityFrameworkCore;
using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Data;

/// <summary>
/// 데이터베이스 컨텍스트
/// </summary>
public class ApplicationDbContext : DbContext
{
    public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
        : base(options)
    {
    }

    public DbSet<User> Users { get; set; }
    public DbSet<Project> Projects { get; set; }
    public DbSet<Circuit> Circuits { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // User 엔티티 설정
        modelBuilder.Entity<User>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.HasIndex(e => e.FirebaseUid).IsUnique();
            entity.HasIndex(e => e.Email);
            entity.Property(e => e.Email).IsRequired();
        });

        // Project 엔티티 설정
        modelBuilder.Entity<Project>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.HasIndex(e => e.OwnerId);
            entity.Property(e => e.Name).IsRequired();
        });

        // Circuit 엔티티 설정
        modelBuilder.Entity<Circuit>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.HasIndex(e => e.ProjectId);
            entity.Property(e => e.Name).IsRequired();
        });
    }
}
