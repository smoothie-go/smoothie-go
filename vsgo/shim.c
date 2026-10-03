#include "shim.h"
#include "_cgo_export.h"

#ifdef _WIN32
#include <windows.h>
#else
#include <dlfcn.h>
#endif

const VSAPI *vsgo_api;

const char *vsgo_load(const char *name) {
    typedef const VSAPI *(VS_CC *getAPI)(int version);
#ifdef _WIN32
    HMODULE lib = LoadLibraryA(name);
    if (!lib)
        return "cannot open the library";
    getAPI get = (getAPI)GetProcAddress(lib, "getVapourSynthAPI");
#else
    void *lib = dlopen(name, RTLD_NOW);
    if (!lib)
        return dlerror();
    getAPI get = (getAPI)dlsym(lib, "getVapourSynthAPI");
#endif
    if (!get)
        return "getVapourSynthAPI not found";
    vsgo_api = get(VAPOURSYNTH_API_VERSION);
    return vsgo_api ? NULL : "api 4 is not supported";
}

static void VS_CC logMessage(int level, const char *msg, void *handle) { vsgoLog(level, (char *)msg, (uintptr_t)handle); }
static void VS_CC logFree(void *handle) { vsgoLogFree((uintptr_t)handle); }

void vsgo_addLogHandler(VSCore *core, uintptr_t handle) {
    vsgo_api->addLogHandler(logMessage, logFree, (void *)handle, core);
}
