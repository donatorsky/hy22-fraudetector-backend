import re
import pickle
from keras.utils import pad_sequences
import keras
from keras import backend as K
import sys

def remove_emojis(data):
    emoj = re.compile("["
        u"\U0001F600-\U0001F64F"  # emoticons
        u"\U0001F300-\U0001F5FF"  # symbols & pictographs
        u"\U0001F680-\U0001F6FF"  # transport & map symbols
        u"\U0001F1E0-\U0001F1FF"  # flags (iOS)
        u"\U00002500-\U00002BEF"  # chinese char
        u"\U00002702-\U000027B0"
        u"\U00002702-\U000027B0"
        u"\U000024C2-\U0001F251"
        u"\U0001f926-\U0001f937"
        u"\U00010000-\U0010ffff"
        u"\u2640-\u2642" 
        u"\u2600-\u2B55"
        u"\u200d"
        u"\u23cf"
        u"\u23e9"
        u"\u231a"
        u"\ufe0f"  # dingbats
        u"\u3030"
                      "]+", re.UNICODE)
    return re.sub(emoj, '', data)

def prepare_data(text):
    text = text.replace('\n', ' ')
    text = text.lower()
    text = remove_emojis(text)
    return text

def recall_m(y_true, y_pred):
    true_positives = K.sum(K.round(K.clip(y_true * y_pred, 0, 1)))
    possible_positives = K.sum(K.round(K.clip(y_true, 0, 1)))
    recall = true_positives / (possible_positives + K.epsilon())
    return recall

def precision_m(y_true, y_pred):
    true_positives = K.sum(K.round(K.clip(y_true * y_pred, 0, 1)))
    predicted_positives = K.sum(K.round(K.clip(y_pred, 0, 1)))
    precision = true_positives / (predicted_positives + K.epsilon())
    return precision

def f1_m(y_true, y_pred):
    precision = precision_m(y_true, y_pred)
    recall = recall_m(y_true, y_pred)
    return 2*((precision*recall)/(precision+recall+K.epsilon()))

# def main():
    # print(sys.argv)
    # text = str(sys.argv[1])
    # text = prepare_data(text)
    
    # Loading Tokenizer
#     with open('tokenizer.pickle', 'rb') as handle:
#         tokenizer = pickle.load(handle)
#     text_tokenized= tokenizer.texts_to_sequences([text])
#     maxlen = 400
#     text_tokenized_final = pad_sequences(text_tokenized, maxlen=maxlen)
        
#     dependencies = {
#         'f1_m': f1_m
#     }
#     # Loading model
#     model = keras.models.load_model('compliance_model', custom_objects=dependencies)
#     result = model.predict(text_tokenized_final)
#     print(result)
#     return result

# def predict():
    
# if __name__ == '__main__':
#     print("asd")
#     main()