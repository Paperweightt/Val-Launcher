#include <windows.h>
#include <shellapi.h>

LRESULT CALLBACK HookedWndProc(HWND hwnd, UINT msg, WPARAM wParam,
                               LPARAM lParam) {
  if (msg == WM_DROPFILES) {
    HDROP hDrop = (HDROP)wParam;
    wchar_t filePath[MAX_PATH];

    if (DragQueryFileW(hDrop, 0, filePath, MAX_PATH)) {
      // Do something with filePath here
      MessageBoxW(NULL, filePath, L"Dropped File", MB_OK);
    }

    DragFinish(hDrop);
  }

  // Call default window procedure
  return DefWindowProcW(hwnd, msg, wParam, lParam);
}

__declspec(dllexport) void SetupDragAndDrop(HWND hwnd) {
  DragAcceptFiles(hwnd, TRUE);
}
