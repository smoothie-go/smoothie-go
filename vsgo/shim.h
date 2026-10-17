#ifndef VSGO_SHIM_H
#define VSGO_SHIM_H

#include "VapourSynth4.h"

extern const VSAPI *vsgo_api;

const char *vsgo_load(const char *name);
void vsgo_addLogHandler(VSCore *core, uintptr_t handle);

// cgo cant call function pointers, so every VSAPI entry in use has a wrapper
static inline VSCore *vsgo_createCore(void) { return vsgo_api->createCore(0); }
static inline void vsgo_freeCore(VSCore *core) { vsgo_api->freeCore(core); }
static inline void vsgo_getCoreInfo(VSCore *core, VSCoreInfo *info) { vsgo_api->getCoreInfo(core, info); }

static inline VSPlugin *vsgo_getPluginByNamespace(const char *ns, VSCore *core) { return vsgo_api->getPluginByNamespace(ns, core); }
static inline VSMap *vsgo_invoke(VSPlugin *plugin, const char *name, const VSMap *args) { return vsgo_api->invoke(plugin, name, args); }

static inline VSMap *vsgo_createMap(void) { return vsgo_api->createMap(); }
static inline void vsgo_freeMap(VSMap *map) { vsgo_api->freeMap(map); }
static inline const char *vsgo_mapGetError(const VSMap *map) { return vsgo_api->mapGetError(map); }
static inline int64_t vsgo_mapGetInt(const VSMap *map, const char *key, int *error) { return vsgo_api->mapGetInt(map, key, 0, error); }
static inline VSNode *vsgo_mapGetNode(const VSMap *map, const char *key) { return vsgo_api->mapGetNode(map, key, 0, NULL); }
static inline void vsgo_mapSetInt(VSMap *map, const char *key, int64_t i) { vsgo_api->mapSetInt(map, key, i, maAppend); }
static inline void vsgo_mapSetFloat(VSMap *map, const char *key, double d) { vsgo_api->mapSetFloat(map, key, d, maAppend); }
static inline void vsgo_mapSetString(VSMap *map, const char *key, const char *s, int size) { vsgo_api->mapSetData(map, key, s, size, dtUtf8, maAppend); }
static inline void vsgo_mapSetNode(VSMap *map, const char *key, VSNode *node) { vsgo_api->mapSetNode(map, key, node, maAppend); }

static inline void vsgo_freeNode(VSNode *node) { vsgo_api->freeNode(node); }
static inline const VSVideoInfo *vsgo_getVideoInfo(VSNode *node) { return vsgo_api->getVideoInfo(node); }
static inline const VSFrame *vsgo_getFrame(int n, VSNode *node, char *err, int errSize) { return vsgo_api->getFrame(n, node, err, errSize); }

static inline void vsgo_freeFrame(const VSFrame *f) { vsgo_api->freeFrame(f); }
static inline const VSMap *vsgo_getFrameProperties(const VSFrame *f) { return vsgo_api->getFramePropertiesRO(f); }
static inline ptrdiff_t vsgo_getStride(const VSFrame *f, int plane) { return vsgo_api->getStride(f, plane); }
static inline const uint8_t *vsgo_getReadPtr(const VSFrame *f, int plane) { return vsgo_api->getReadPtr(f, plane); }
static inline int vsgo_getFrameWidth(const VSFrame *f, int plane) { return vsgo_api->getFrameWidth(f, plane); }
static inline int vsgo_getFrameHeight(const VSFrame *f, int plane) { return vsgo_api->getFrameHeight(f, plane); }

#endif
