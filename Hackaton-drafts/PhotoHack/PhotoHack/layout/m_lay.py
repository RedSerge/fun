from PIL import Image, ImageDraw
import numpy as np
import math
import os, shutil


def crop(infile_r, width, height):
    im = Image.open(infile_r)
    imgwidth, imgheight = im.size
    for i in range(imgheight//height):
        for j in range(imgwidth//width):
            box = (j*width, i*height, (j+1)*width, (i+1)*height)
            yield im.crop(box)


def crop_mulpty(infile_m, poly):
    im = Image.open(infile_m).convert("RGBA")
    imArray = np.asarray(im)
    maskIm = Image.new("L", (imArray.shape[1], imArray.shape[0]), 0)
    ImageDraw.Draw(maskIm).polygon(poly, outline=1, fill=1)
    mask = np.array(maskIm)
    newImArray = np.empty(imArray.shape,dtype='uint8')
    newImArray[:,:,:3] = imArray[:,:,:3]
    newImArray[:,:,3]=mask*255
    img_t = Image.fromarray(newImArray, "RGBA")
    return img_t


def basic():
    name_file=input("имя файла № 1:  ")
    Image.open('..\data\in\\' + str(name_file)).convert("RGBA").save('..\data\out\opt\АВА.png')
    infile='..\data\out\opt\АВА.png'
    with Image.open(infile) as img:
        width, height = img.size
        min_s=min(width, height)
        img_r=img.resize((min_s,min_s))
        img_r.save('..\data\out\opt\АВА_r.png')
        infile_r='..\data\out\opt\АВА_r.png'
    count=input("Количество осколков:  ")
    count=int(math.sqrt(int(count)/5))
    h_size=int(min_s/count)
    w_size=int(min_s/count)
    for k, piece in enumerate(crop(infile_r, w_size, h_size)):
        img_r=Image.new("RGBA", (w_size, h_size), 255)
        img_r.paste(piece)
        path=os.path.join('..\data\out\squares',"IMG-{}.png".format(k+1))
        img_r.save(path)

        img_1 = crop_mulpty(path, [(0,0),(0,int(h_size/2)),(int(w_size/2),0)])
        path_1=os.path.join('../js/NewFl2/NewFl/static/debris/',"IMG-{}-{}.png".format(k,1))
        img_1.save(path_1)

        img_2 = crop_mulpty(path,[(0,int(h_size/2)),(0,int(h_size)),(int(w_size/2),0)])
        path_2=os.path.join('../js/NewFl2/NewFl/static/debris/',"IMG-{}-{}.png".format(k,2))
        img_2.save(path_2)

        img_3 = crop_mulpty(path,[(0,int(h_size)),(int(w_size/2),int(h_size/2)),(int(w_size/2),0)])
        path_3=os.path.join('../js/NewFl2/NewFl/static/debris/',"IMG-{}-{}.png".format(k,3))
        img_3.save(path_3)

        img_4 = crop_mulpty(path,[(0,int(h_size)),(int(w_size),int(h_size)),(int(w_size/2),int(h_size/2))])
        path_4=os.path.join('../js/NewFl2/NewFl/static/debris/',"IMG-{}-{}.png".format(k,4))
        img_4.save(path_4)

        img_5 = crop_mulpty(path,[(int(w_size/2),0),(int(w_size/2),int(h_size/2)),(int(w_size),int(h_size)),(int(w_size),0)])
        path_5=os.path.join('../js/NewFl2/NewFl/static/debris/',"IMG-{}-{}.png".format(k,5))
        img_5.save(path_5)





    name_file=input("имя файла № 2:  ")
    Image.open('..\data\in\\' + str(name_file)).convert("RGBA").save('..\data\out\opt\АВА1.png')
    infile='..\data\out\opt\АВА1.png'
    with Image.open(infile) as img:
        width, height = img.size
        min_s=min(width, height)
        img_r=img.resize((min_s,min_s))
        img_r.save('..\data\out\opt\АВА1_r.png')
        infile_r='..\data\out\opt\АВА1_r.png'
    h_size=int(min_s/count)
    w_size=int(min_s/count)
    for k, piece in enumerate(crop(infile_r, w_size, h_size)):
        img_r=Image.new("RGBA", (w_size, h_size), 255)
        img_r.paste(piece)
        path=os.path.join('..\data\out\squares_two',"IMG-{}.png".format(k+1))
        img_r.save(path)

        img_1 = crop_mulpty(path, [(0,0),(0,int(h_size/2)),(int(w_size/2),0)])
        path_1=os.path.join('../js/NewFl2/NewFl/static/debris_two/',"IMG-{}-{}.png".format(k,1))
        img_1.save(path_1)

        img_2 = crop_mulpty(path,[(0,int(h_size/2)),(0,int(h_size)),(int(w_size/2),0)])
        path_2=os.path.join('../js/NewFl2/NewFl/static/debris_two/',"IMG-{}-{}.png".format(k,2))
        img_2.save(path_2)

        img_3 = crop_mulpty(path,[(0,int(h_size)),(int(w_size/2),int(h_size/2)),(int(w_size/2),0)])
        path_3=os.path.join('../js/NewFl2/NewFl/static/debris_two/',"IMG-{}-{}.png".format(k,3))
        img_3.save(path_3)

        img_4 = crop_mulpty(path,[(0,int(h_size)),(int(w_size),int(h_size)),(int(w_size/2),int(h_size/2))])
        path_4=os.path.join('../js/NewFl2/NewFl/static/debris_two/',"IMG-{}-{}.png".format(k,4))
        img_4.save(path_4)

        img_5 = crop_mulpty(path,[(int(w_size/2),0),(int(w_size/2),int(h_size/2)),(int(w_size),int(h_size)),(int(w_size),0)])
        path_5=os.path.join('../js/NewFl2/NewFl/static/debris_two/',"IMG-{}-{}.png".format(k,5))
        img_5.save(path_5)


def re_move_files():
    folder = ['../data/out/opt', '../data/out/squares', '../data/out/squares_two', '../js/NewFl2/NewFl/static/debris','../js/NewFl2/NewFl/static/debris_two']
    for j in range(len(folder)-1):
        for file in os.listdir(folder[j]):
            os.remove(folder[j]+'/'+file)



if __name__=='__main__':
    re_move_files()
    basic()
