import cv2
import os
import time

import numpy as np


def reader(img_bytes):
    img_bytes = img_bytes.read()
    nparr = np.frombuffer(img_bytes, np.uint8)
    imagem = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    cinza = cv2.cvtColor(imagem, cv2.COLOR_BGR2GRAY)
    classificador = cv2.CascadeClassifier(cv2.data.haarcascades + "haarcascade_frontalface_default.xml")
    faces = classificador.detectMultiScale(cinza, scaleFactor=1.1, minNeighbors=5, minSize=(30, 30))

    for (x, y, w, h) in faces:
        cv2.rectangle(imagem, (x, y), (x + w, y + h), (0, 255, 0), 2)

    path = r'/home/felipe/Aplicação_UVA_Paralela/PyServer/resultados'
    os.chdir(path)

    status, buffer = cv2.imencode('.png', imagem)
    img_name = 'result.png'
    with open(img_name, 'wb') as f:
        f.write(buffer)

    if status:
        print("Imagem salva com sucesso")
    else:
        print("O salvamento da imagem falhou")
    cv2.waitKey(0)
    cv2.destroyAllWindows()
