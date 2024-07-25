const std = @import("std");
const sdl = @import("SDL.zig/build.zig");

pub fn build(b: *std.Build) !void {
    // Determine compilation target
    const target = b.standardTargetOptions(.{});

    // Create a new instance of the SDL2 Sdk
    const sdk = sdl.init(b, null, null);

    // Create executable for our example
    const demo_basic = b.addExecutable(.{
        .name = "tetris",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
    });

    sdk.link(demo_basic, .dynamic, sdl.Library.SDL2); // link SDL2 as a shared library

    // Add "sdl2" package that exposes the SDL2 api (like SDL_Init or SDL_CreateWindow)
    demo_basic.root_module.addImport("sdl2", sdk.getNativeModule());

    // Install the executable into the prefix when invoking "zig build"
    b.installArtifact(demo_basic);
}
