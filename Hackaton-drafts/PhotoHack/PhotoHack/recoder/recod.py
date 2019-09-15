import numpy as np
import cv2
import pyscreenshot as ImageGrab

def freeze_support():
    fourcc = cv2.VideoWriter_fourcc(*'XVID')
    out = cv2.VideoWriter("record.avi", fourcc, 5.0, (1366, 768))

    i=0
    while True:
        img = ImageGrab.grab()
        img_np = np.array(img)
        frame = cv2.cvtColor(img_np, cv2.COLOR_BGR2RGB)
	#cv2.imshow("Screen", frame)
        out.write(frame)

        k = cv2.waitKey(33)
        if k == ord('s'):
            image = pyautogui.screenshot(region=(0,0, 768, 1366))
            cv2.imwrite('scrin_{}.png'.format(i), image)
            i+=1

        if k == 27:
            break

    out.release()
    cv2.destroyAllWindows()

if __name__ == '__main__':
    freeze_support()
